package controller

import (
	"net/http"

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
	room := ""
	exist := true
	var err error

	for exist {
		room = ctrl.svc.GenerateRoom(6)
		exist, err = ctrl.svc.CheckRoomExist(room)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	ctx.JSON(http.StatusOK, gin.H{"room": room})
}
