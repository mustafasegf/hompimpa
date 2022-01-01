package room

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"github.com/mustafasegf/hompimpa/entity"
)

type Service struct {
	repo *Repo
}

func NewService(repo *Repo) *Service {
	return &Service{
		repo: repo,
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func (r *Service) GenerateRoom(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func (r *Service) CheckRoomExist(room string) (exist bool, err error) {
	ctx := context.Background()
	channels, err := r.repo.Pub.PubSubChannels(ctx, room).Result()
	if err != nil {
		return
	}
	for _, channel := range channels {
		if channel == room {
			exist = true
			return
		}
	}
	exist = false
	return
}

func (r *Service) SubscribeToRoom(room string) (sub *redis.PubSub) {
	ctx := context.Background()
	sub = r.repo.Sub.Subscribe(ctx, room)
	return
}

func (r *Service) PublishToRoom(room, data string) (err error) {
	ctx := context.Background()
	err = r.repo.Pub.Publish(ctx, room, data).Err()
	return
}

func (r *Service) GetRoomData(room, data string) (roomData *entity.RoomData, err error) {
	ctx := context.Background()
	res := r.repo.Pub.Get(ctx, room)
	if res.Err() != nil {
		err = fmt.Errorf("get: %v", err)
		return
	}

	raw, err := res.Bytes()
	if err != nil {
		err = fmt.Errorf("byte: %v", err)
		return
	}

	roomData = &entity.RoomData{}
	err = json.Unmarshal(raw, roomData)
	if err != nil {
		err = fmt.Errorf("unmarshall: %v", err)
		roomData = nil
		return
	}

	return
}

func (r *Service) CreateRoomData(room, name string) (roomData *entity.RoomData, err error) {
	users := map[string]entity.User{
		name: {
			Name: name,
		},
	}

	roomData = &entity.RoomData{
		Status: "waiting",
		Owner:  name,
		Users:  users,
	}

	roomString, err := json.Marshal(roomData)
	if err != nil {
		err = fmt.Errorf("marshal: %v", err)
		roomData = nil
		return
	}

	ctx := context.Background()
	r.repo.Pub.Set(ctx, room, roomString, 0)
	return
}

func (r *Service) AddRoomData(room string, user entity.User, oldRoom *entity.RoomData) (roomData *entity.RoomData, err error) {
	users := oldRoom.Users
	users[user.Name] = user
	oldRoom.Users = users

	if len(users) > 2 && r.CalculateGame(users) {
		oldRoom.Status = "all"
	}

	roomString, err := json.Marshal(oldRoom)
	if err != nil {
		err = fmt.Errorf("marshal: %v", err)
		return
	}

	ctx := context.Background()
	r.repo.Pub.Set(ctx, room, roomString, 0)
	roomData = oldRoom
	return
}

func (r *Service) RemoveUser(room, name string) (err error) {
	roomData, err := r.GetRoomData(room, name)
	if err != nil {
		err = fmt.Errorf("remove: %v", err)
		return
	}
	users := roomData.Users
	delete(users, name)
	roomData.Users = users

	if len(users) > 2 && r.CalculateGame(users) {
		roomData.Status = "all"
	}

	ctx := context.Background()
	if len(users) == 0 {
		err = r.repo.Pub.Del(ctx, room).Err()
		return
	}

	roomString, err := json.Marshal(roomData)
	if err != nil {
		err = fmt.Errorf("marshal: %v", err)
		return
	}

	r.repo.Pub.Set(ctx, room, roomString, 0)
	return
}

func (r *Service) CalculateGame(users map[string]entity.User) (valid bool) {
	for _, user := range users {
		if user.Hand == nil {
			valid = false
			return
		}
	}
	valid = true
	return
}

func (r *Service) ReadMessage(ctx context.Context, ws *websocket.Conn, room string) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			_, msg, err := ws.ReadMessage()
			if err != nil {
				return
			}
			userData := entity.User{}
			json.Unmarshal(msg, &userData)

			// get room
			roomData, err := r.GetRoomData(room, string(msg))
			if err != nil {
				log.Println("get room:", err)
			}

			// add data
			if roomData != nil {
				roomData, err = r.AddRoomData(room, userData, roomData)
				if err != nil {
					log.Println("add room data:", err)
				}
			} else {
				roomData, err = r.CreateRoomData(room, userData.Name)
				if err != nil {
					log.Println("create room data:", err)
				}
			}

			roomByte, err := json.Marshal(roomData)
			// send room info
			r.PublishToRoom(room, string(roomByte))
		}
	}
}

func (r *Service) WriteMessage(ctx context.Context, ws *websocket.Conn, chn <-chan *redis.Message) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-chn:
			ws.WriteJSON(msg)
		}
	}
}

func (r *Service) Ping(ctx context.Context, cancel context.CancelFunc, ws *websocket.Conn, ticker *time.Ticker, room, name string) {
	defer func() {
		r.RemoveUser(room, name)
		fmt.Printf("%s disconected\n", name)
		cancel()
	}()
	for {
		select {
		case <-ticker.C:
			ws.SetWriteDeadline(time.Now().Add(WriteWait))
			if err := ws.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
