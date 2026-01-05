package routes

import (
	"server/controllers"

	"github.com/gin-gonic/gin"
)

func SetUpRoutes(router *gin.Engine) {

	router.POST("/login", controllers.LoginControllers)
	// router.POST("/register", controllers.RegisterController)
	router.GET("/auth/:provider/callback", controllers.GetAuthCallBackFunction)
	router.GET("/auth/:provider", controllers.GoogleLogin)
	// router.DELETE("/logout", controllers.Logout)
}
