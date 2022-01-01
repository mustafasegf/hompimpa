package room

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Controller struct {
	svc  *Service
	upgr websocket.Upgrader
}

func NewController(svc *Service, upgr websocket.Upgrader) *Controller {
	return &Controller{
		svc:  svc,
		upgr: upgr,
	}
}

func (ctrl *Controller) CreateRoom(ctx *gin.Context) {
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

func (ctrl *Controller) Connect(ctx *gin.Context) {
	room := ctx.Param("room")
	if len(room) != 6 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "wrong length"})
		return
	}

	ws, err := ctrl.upgr.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	sub := ctrl.svc.SubscribeToRoom(room)
	chn := sub.Channel()

	ticker := time.NewTicker(PingPeriod)
	c := context.Background()
	c, cancel := context.WithCancel(c)
	defer func() {
		ticker.Stop()
		ws.Close()
		sub.Close()
	}()

	_, name, _ := ws.ReadMessage()

	go ctrl.svc.ReadMessage(c, ws, room)
	go ctrl.svc.WriteMessage(c, ws, chn)
	go ctrl.svc.Ping(c, cancel, ws, ticker, room, string(name))

	for {
		select {
		case <-c.Done():
			return
		}
	}
}
