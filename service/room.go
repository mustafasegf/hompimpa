package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"github.com/mustafasegf/hompimpa/constant"
	"github.com/mustafasegf/hompimpa/entity"
	"github.com/mustafasegf/hompimpa/repository"
)

type Room struct {
	repo *repository.Room
}

func NewRoomService(repo *repository.Room) *Room {
	return &Room{
		repo: repo,
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func (r *Room) GenerateRoom(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func (r *Room) CheckRoomExist(room string) (exist bool, err error) {
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

func (r *Room) SubscribeToRoom(room string) (sub *redis.PubSub) {
	ctx := context.Background()
	sub = r.repo.Sub.Subscribe(ctx, room)
	return
}

func (r *Room) PublishToRoom(room, data string) (err error) {
	ctx := context.Background()
	err = r.repo.Pub.Publish(ctx, room, data).Err()
	return
}

func (r *Room) GetRoomData(room, data string) (roomData *entity.RoomData, err error) {
	ctx := context.Background()
	res := r.repo.Pub.Get(ctx, room)
	if res.Err() != nil {
		err = fmt.Errorf("get: %v", err)
		return
	}
	err = nil
	fmt.Printf("res: %s\n", res.String())
	raw, err := res.Bytes()
	if err != nil {
		err = fmt.Errorf("byte: %v", err)
		return
	}
	fmt.Printf("raw: %s\n", raw)

	roomData = &entity.RoomData{}
	err = json.Unmarshal(raw, roomData)
	if err != nil {
		err = fmt.Errorf("unmarshall: %v", err)
		roomData = nil
		return
	}
	fmt.Printf("roomData: %#v\n", roomData)
	return
}

func (r *Room) CreateRoomData(room, name string) (roomData *entity.RoomData, err error) {
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

func (r *Room) AddRoomData(room, name string, oldRoom *entity.RoomData) (roomData *entity.RoomData, err error) {
	users := oldRoom.Users
	users[name] = entity.User{Name: name}
	oldRoom.Users = users

	roomString, err := json.Marshal(oldRoom)
	if err != nil {
		err = fmt.Errorf("marshal: %v", err)
		return
	}

	fmt.Printf("room string: %s err:%v\n", roomString, err)
	ctx := context.Background()
	r.repo.Pub.Set(ctx, room, roomString, 0)
	roomData = oldRoom
	return
}

func (r *Room) ReadMessage(ctx context.Context, ws *websocket.Conn, room string) {
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

			if roomData != nil {
				roomData, err = r.AddRoomData(room, userData.Name, roomData)
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

func (r *Room) WriteMessage(ctx context.Context, ws *websocket.Conn, chn <-chan *redis.Message) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-chn:
			ws.WriteJSON(msg)
		}
	}
}

func (r *Room) Ping(ctx context.Context, cancel context.CancelFunc, ws *websocket.Conn, ticker *time.Ticker) {
	defer cancel()
	for {
		select {
		case <-ticker.C:
			ws.SetWriteDeadline(time.Now().Add(constant.WriteWait))
			if err := ws.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
