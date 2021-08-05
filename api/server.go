package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"github.com/mustafasegf/hompimpa/util"
)

type Server struct {
	config util.Config
	router *gin.Engine
	pub    *redis.Client
	sub    *redis.Client
	upgr   websocket.Upgrader
}

func MakeServer(config util.Config, pub *redis.Client, sub *redis.Client) Server {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	router := gin.Default()
	server := Server{
		config: config,
		router: router,
		pub:    pub,
		sub:    sub,
		upgr:   upgrader,
	}
	return server
}

func (s *Server) RunServer() {
	s.setupRouter()
	serverString := fmt.Sprintf(":%s", s.config.ServerPort)
	s.router.Run(serverString)
}
