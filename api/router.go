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
	roomRepo := repository.NewRoomRepo(s.pub, s.sub)
	roomSvc := service.NewRoomService(roomRepo)
	roomCtlr := controller.NewRoomController(roomSvc, s.upgr)

	root := s.router.Group("/")
	root.GET(":room", roomCtlr.Connect)

	api := s.router.Group("/api")
	room := api.Group("/room")

	room.GET("create", roomCtlr.CreateRoom)

}
