package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/mustafasegf/hompimpa/service"
)

type Room struct {
	svc *service.Room
}

func NewRoomController(svc *service.Room) *Room {
	return &Room{
		svc: svc,
	}
}

func (ctrl *Room) CreateRoom(ctx *gin.Context) {
	
}
