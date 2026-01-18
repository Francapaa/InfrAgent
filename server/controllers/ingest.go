package controllers

import (
	"net/http"
	models "server/model"
	"server/service"
	"strings"

	"github.com/gin-gonic/gin"
)

type IngestHandlerController struct {
	ingestHandler service.IngestHandler
}

func NewNewEventInRequest(s service.IngestHandler) *IngestHandlerController {
	return &IngestHandlerController{ingestHandler: s}
}

func (IHC *IngestHandlerController) NewEventInRequest(c *gin.Context) {

	authHeader := c.GetHeader("Authorization")

	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Se requiere API KEY"})
		return
	}
	apiKey := strings.TrimPrefix(authHeader, "Bearer")

	var event models.Event

	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "formato invalido para evento"})
		return
	}

	err := IHC.ingestHandler.NewEventInRequestService(c.Request.Context(), apiKey, &event)

	if err != nil {
		if err.Error() == "unauthorized" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "faltan credenciales para poder validar la peticion"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error interno, intentelo mas tarde"})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"status": "event_queued"})

}
