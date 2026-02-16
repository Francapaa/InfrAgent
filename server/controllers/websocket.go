package controllers

import (
	"net/http"
	"server/ws"

	"github.com/gin-gonic/gin"
)

type WebSocketController struct{}

func (wsc *WebSocketController) HandleWebSocket(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user ID"})
		return
	}

	hub := ws.GetHub()
	ws.ServeWs(hub, ctx.Writer, ctx.Request, userIDStr)
}
