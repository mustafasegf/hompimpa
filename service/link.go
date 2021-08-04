package service

import (
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
