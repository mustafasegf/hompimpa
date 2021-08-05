package service

import (
	"context"
	"math/rand"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"github.com/mustafasegf/hompimpa/constant"
	"github.com/mustafasegf/hompimpa/repository"
)

type WsData struct {
	Name string
	Data interface{}
}

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

func (r *Room) ReadMessage(ctx context.Context, ws *websocket.Conn, room string) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			_, message, err := ws.ReadMessage()
			if err != nil {
				return
			}
			r.PublishToRoom(room, string(message))
		}
	}
}

func (r *Room) WriteMessage(ctx context.Context, ws *websocket.Conn, chn <-chan *redis.Message) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-chn:
			ws.WriteJSON(WsData{Data: msg.String()})
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
