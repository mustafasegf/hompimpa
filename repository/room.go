package repository

import (
	"github.com/go-redis/redis/v8"
)

type Room struct {
	Rdb *redis.Client
}

func NewRoomRepo(rdb *redis.Client) *Room {
	return &Room{
		Rdb: rdb,
	}
}
