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

	ticker := time.NewTicker(constant.PingPeriod)
	c := context.Background()
	c, cancel := context.WithCancel(c)
	defer func() {
		ticker.Stop()
		ws.Close()
		sub.Close()
	}()

	go ctrl.svc.ReadMessage(c, ws, room)
	go ctrl.svc.WriteMessage(c, ws, chn)
	go ctrl.svc.Ping(c, cancel, ws, ticker)

	for {
		select {
		case <-c.Done():
			return
		}
	}
}

/* client 1 connect to ws
send {name: test1}
saved to rejson
get data from rejson
send {
	status: waiting,
	owner: test1,
	users: [
		{
			name: test1,
			hand: null
		}
	]
}

client 2 connected
send {name: test2}
saved to rejson
get data from rejson
send {
	status: waiting,
	leader: test1,
	users: [
		{
			name: test1,
			hand: null
		},
		{
			name: test2,
			hand: null
		}
	]
}
*/
