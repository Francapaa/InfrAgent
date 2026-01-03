package routes

import (
	"server/controllers"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func setUpRoutes(router *gin.Engine, db *mongo.Database) {

	router.POST("/login", controllers.LoginControllers)
	// router.POST("/register", controllers.RegisterController)
	router.GET("/auth/callback", controllers.GetAuthCallBackFunction)
	router.GET("/auth/google", controllers.GoogleLogin)
	// router.DELETE("/logout", controllers.Logout)
}
