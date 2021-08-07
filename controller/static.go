package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Static struct {
}

func NewStaticController() *Static {
	return &Static{}
}

func (server *Static) Index(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "index.html", "")
}

func (server *Static) RenderRoom(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "room.html", gin.H{"id": ctx.Param("room")})
}
