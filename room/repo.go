package room

import (
	"github.com/go-redis/redis/v8"
)

type Repo struct {
	Pub *redis.Client
	Sub *redis.Client
}

func NewRepo(pub *redis.Client, sub *redis.Client) *Repo {
	return &Repo{
		Pub: pub,
		Sub: sub,
	}
}
