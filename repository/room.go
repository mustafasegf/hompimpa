package repository

import (
	"github.com/go-redis/redis/v8"
)

type Room struct {
	Pub *redis.Client
	Sub *redis.Client
}

func NewRoomRepo(pub *redis.Client, sub *redis.Client) *Room {
	return &Room{
		Pub: pub,
		Sub: sub,
	}
}
