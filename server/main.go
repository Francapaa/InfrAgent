package main

import (
	"log"
	"server/controllers"
	"server/internal/database"
	"server/repositories"
	"server/routes"
	"server/service"
	"server/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error cargando .env")
	}

	db := database.ConnectDatabase()
	defer db.Close()
	sqlDB := database.GetSQLDB()
	utils.NewAuth()

	loginService := service.NewLogin(repositories.NewPostgresStorage(sqlDB))

	loginController := controllers.NewLoginController(loginService)

	wsController := controllers.NewWebSocketController()

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	setupRoutes := routes.NewSetUpRoutes(loginController, wsController)

	setupRoutes.SetUpRoutes(router)

	log.Println("🚀 Servidor corriendo en http://localhost:8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Error al iniciar el servidor:", err)
	}
}
