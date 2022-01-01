package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mustafasegf/hompimpa/room"
)

type Route struct {
	router *gin.Engine
}

func (s *Server) setupRouter() {
	s.router.LoadHTMLGlob("templates/*")
	s.router.Static("/static", "./static")

	roomRepo := room.NewRepo(s.pub, s.sub)
	roomSvc := room.NewService(roomRepo)
	roomCtlr := room.NewController(roomSvc, s.upgr)

	root := s.router.Group("/")
	root.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", "")
	})

	root.GET(":room", func(ctx *gin.Context) {
		room := ctx.Param("room")
		if len(room) != 6 {
			ctx.Redirect(http.StatusTemporaryRedirect, "/")
			return
		}
		ctx.HTML(http.StatusOK, "room.html", gin.H{"id": room})
	})

	api := s.router.Group("/api/room")

	api.GET("create", roomCtlr.CreateRoom)
	api.GET(":room", roomCtlr.Connect)

}
