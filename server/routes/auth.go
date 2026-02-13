package routes

import (
	"server/controllers"
	"server/middleware"

	"github.com/gin-gonic/gin"
)

type SetUpRoutes struct {
	controllers     *controllers.LoginController
	wsController    *controllers.WebSocketController
	agentController *controllers.AgentController
}

func (sp *SetUpRoutes) SetUpRoutes(router *gin.Engine) {
	router.GET("/auth/:provider/callback", sp.controllers.GetAuthCallBackFunction)
	router.GET("/auth/:provider", sp.controllers.GoogleLogin)
	router.GET("/ws", sp.wsController.HandleWebSocket)

	// Rutas protegidas con JWT
	authorized := router.Group("/auth")
	authorized.Use(middleware.JWTMiddleware())
	{
		authorized.GET("/me", sp.controllers.GetCurrentUser)
		authorized.POST("/complete-registration", sp.controllers.CompleteRegistration)
	}

	// Rutas del agente (protegidas con JWT)
	agentRoutes := router.Group("/api/agent")
	agentRoutes.Use(middleware.JWTMiddleware())
	{
		agentRoutes.GET("/state", sp.agentController.GetAgentState)
		agentRoutes.GET("/actions", sp.agentController.GetLastRecentActions)
	}
}

func NewSetUpRoutes(loginController *controllers.LoginController, wsController *controllers.WebSocketController, agentController *controllers.AgentController) *SetUpRoutes {
	return &SetUpRoutes{
		controllers:     loginController,
		wsController:    wsController,
		agentController: agentController,
	}
}
