package service

import (
	"context"
	"math/rand"
	"time"

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
	channels, err := r.repo.Rdb.PubSubChannels(ctx, room).Result()
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
