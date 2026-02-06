package controllers

import (
	"server/ws"

	"github.com/gin-gonic/gin"
)

type WebSocketController struct{}

func (wsc *WebSocketController) HandleWebSocket(c *gin.Context) {
	hub := ws.GetHub()
	ws.ServeWs(hub, c.Writer, c.Request)
}
