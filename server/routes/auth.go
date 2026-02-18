package routes

import (
	"server/controllers"
	"server/middleware"
	"server/repositories"

	"github.com/gin-gonic/gin"
)

type SetUpRoutes struct {
	controllers      *controllers.LoginController
	wsController     *controllers.WebSocketController
	agentController  *controllers.AgentController
	clientStorage    *repositories.ClientStorage
	ingestController *controllers.IngestHandlerController
}

func (sp *SetUpRoutes) SetUpRoutes(router *gin.Engine) {
	router.GET("/auth/:provider/callback", sp.controllers.GetAuthCallBackFunction)
	router.GET("/auth/:provider", sp.controllers.GoogleLogin)

	// Rutas protegidas con JWT
	authorized := router.Group("/auth")
	authorized.Use(middleware.JWTMiddleware())
	{
		authorized.GET("/me", sp.controllers.GetCurrentUser)
		complete := authorized.Group("/")
		complete.Use(middleware.ProfileCompleteMiddleware(*sp.clientStorage))
		{
			authorized.POST("/complete-registration", sp.controllers.CompleteRegistration)
		}
	}

	// Rutas del agente (protegidas con JWT)
	agentRoutes := router.Group("/api/agent")
	agentRoutes.Use(middleware.JWTMiddleware())
	{
		agentRoutes.GET("/state", sp.agentController.GetAgentState)
		agentRoutes.GET("/actions", sp.agentController.GetLastRecentActions)
	}
	websocketsRoutes := router.Group("/ws")
	websocketsRoutes.Use(middleware.JWTMiddleware())
	{
		router.GET("/", sp.wsController.HandleWebSocket) // ahora podemos identificar segun el cliente
	}
	SDKRoutes := router.Group("/sdk")
	SDKRoutes.Use(middleware.AuthMiddlewareApiKey())
	{
		SDKRoutes.POST("/event", sp.ingestController.NewEventInRequest)
	}
}

func NewSetUpRoutes(loginController *controllers.LoginController, wsController *controllers.WebSocketController, agentController *controllers.AgentController) *SetUpRoutes {
	return &SetUpRoutes{
		controllers:     loginController,
		wsController:    wsController,
		agentController: agentController,
	}
}
