package controller

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/mustafasegf/hompimpa/constant"
	"github.com/mustafasegf/hompimpa/service"
)

type WsData struct {
	Name string
	Data interface{}
}

type Room struct {
	svc  *service.Room
	upgr websocket.Upgrader
}

func NewRoomController(svc *service.Room, upgr websocket.Upgrader) *Room {
	return &Room{
		svc:  svc,
		upgr: upgr,
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

func (ctrl *Room) Connect(ctx *gin.Context) {
	room := ctx.Param("room")
	if len(room) != 6 {
		ctx.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	ws, err := ctrl.upgr.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	sub := ctrl.svc.SubscribeToRoom(room)
	chn := sub.Channel()

	ticker := time.NewTicker(constant.PingPeriod)
	c := context.Background()
	c, cancel := context.WithCancel(c)

	defer func() {
		cancel()
		ticker.Stop()
		ws.Close()
		sub.Close()
	}()

	go func(c context.Context) {
		for {
			select {
			case <-c.Done():
				return
			default:
				_, message, err := ws.ReadMessage()
				if err != nil {
					return
				}
				ctrl.svc.PublishToRoom(room, string(message))
			}
		}
	}(c)

	for {
		select {
		case msg := <-chn:
			ws.WriteJSON(WsData{Data: msg.String()})
		case <-ticker.C:
			ws.SetWriteDeadline(time.Now().Add(constant.WriteWait))
			if err := ws.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
