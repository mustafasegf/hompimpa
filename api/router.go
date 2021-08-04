package api

import (
	"github.com/gin-gonic/gin"
	"github.com/mustafasegf/hompimpa/controller"
	"github.com/mustafasegf/hompimpa/repository"
	"github.com/mustafasegf/hompimpa/service"
)

type Route struct {
	router *gin.Engine
}

func (s *Server) setupRouter() {
	roomRepo := repository.NewRoomRepo(s.rdb)
	roomSvc := service.NewRoomService(roomRepo)
	roomCtlr := controller.NewRoomController(roomSvc)

	api := s.router.Group("/api")
	room := api.Group("/room")

	room.GET("create", roomCtlr.CreateRoom)

}
