package routes

import (
	"server/controllers"

	"github.com/gin-gonic/gin"
)

type SetUpRoutes struct {
	controllers controllers.LoginController
}

func (sp *SetUpRoutes) SetUpRoutes(router *gin.Engine) {

	router.POST("/login", sp.controllers.LoginControllers)
	router.POST("/register", sp.controllers.LocalRegister)
	router.GET("/auth/:provider/callback", sp.controllers.GetAuthCallBackFunction)
	router.GET("/auth/:provider", sp.controllers.GoogleLogin)
	// router.DELETE("/logout", controllers.Logout)
}
