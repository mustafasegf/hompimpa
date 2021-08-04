package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/mustafasegf/hompimpa/util"
)

type Server struct {
	config util.Config
	router *gin.Engine
	rdb    *redis.Client
}

func MakeServer(config util.Config, rdb *redis.Client) Server {
	router := gin.Default()
	server := Server{
		config: config,
		router: router,
		rdb:    rdb,
	}
	return server
}

func (s *Server) RunServer() {
	s.setupRouter()
	serverString := fmt.Sprintf(":%s", s.config.ServerPort)
	s.router.Run(serverString)
}
