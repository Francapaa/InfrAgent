package controllers

import (
	"net/http"
	"server/middleware"
	"server/service"

	"github.com/gin-gonic/gin"
)

// AgentController maneja los endpoints relacionados con el agente
type AgentController struct {
	agentDataService *service.AgentDataService
}

// NewAgentController crea un nuevo controlador de agente
func NewAgentController(agentDataService *service.AgentDataService) *AgentController {
	return &AgentController{
		agentDataService: agentDataService,
	}
}

// GetAgentState devuelve el estado completo del agente del usuario logueado
func (c *AgentController) GetAgentState(ctx *gin.Context) {
	// Obtener el userID del contexto (seteado por el middleware JWT)
	userID, exists := ctx.Get(middleware.UserIDKey)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}

	clientID, ok := userID.(string)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener ID del usuario"})
		return
	}

	// Usar el servicio para obtener el estado del agente (lógica de negocio)
	state, err := c.agentDataService.GetAgentStateForClient(ctx, clientID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, state)
}

func (c *AgentController) GetLastRecentActions(ctx *gin.Context) {

	userID, exists := ctx.Get(middleware.UserIDKey)

	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}

	clientID, ok := userID.(string)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener ID del usuario"})
		return
	}

	actions, err := c.agentDataService.GetLast30ActionsByAgent(ctx, clientID)

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, actions)
	// LLAMAR A SERVICE QUE RETORNE LAS ULTIMA 30 ACCIONES QUE HAYA HECHO EL AGENTE
}
