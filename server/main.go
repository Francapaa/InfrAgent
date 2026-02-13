package main

import (
	"context"
	"log"
	"os"
	"server/controllers"
	"server/internal/database"
	"server/repositories"
	"server/routes"
	"server/service"
	"server/service/agent/llm"
	"server/utils"
	"time"

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

	// Inicializar repositorios
	repo := repositories.NewPostgresStorage(sqlDB)

	// Inicializar cliente Gemini
	geminiAPIKey := os.Getenv("KEY_API_GEMINI")
	if geminiAPIKey == "" {
		log.Fatal("GEMINI_API_KEY no está configurada en las variables de entorno")
	}
	geminiClient := llm.ConnectionToGeminiLLM(geminiAPIKey, "gemini-2.0-flash-exp")

	// Inicializar Agent Engine
	agentEngine := &service.AgentEngine{
		Gemini:  geminiClient,
		Events:  repo,
		Actions: repo,
		Agents:  repo,
		Client:  repo,
	}

	// Crear servicio de datos del agente (lógica de negocio)
	agentDataService := service.NewAgentDataService(repo, agentEngine)

	// Crear controlador del agente (solo HTTP)
	agentController := controllers.NewAgentController(agentDataService)
	startAgentScheduler(agentEngine, repo)
	router := gin.New()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	setupRoutes := routes.NewSetUpRoutes(loginController, wsController, agentController)

	setupRoutes.SetUpRoutes(router)

	log.Println("Servidor corriendo en http://localhost:8080")
	log.Println("Agente autónomo iniciado (tick cada 30 segundos)")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Error al iniciar el servidor:", err)
	}
}

// startAgentScheduler inicia el loop del agente que ejecuta ticks cada 30 segundos
func startAgentScheduler(engine *service.AgentEngine, repo repositories.AgentStorage) {
	ticker := time.NewTicker(30 * time.Second)

	go func() {
		// Ejecutar inmediatamente al iniciar
		log.Println("[Scheduler] Ejecutando tick inicial...")
		runAgentTick(engine, repo)

		// Luego ejecutar cada 30 segundos
		for range ticker.C {
			log.Println("[Scheduler] Ejecutando tick programado...")
			runAgentTick(engine, repo)
		}
	}()
}

// runAgentTick ejecuta un tick para todos los agentes activos
func runAgentTick(engine *service.AgentEngine, repo repositories.AgentStorage) {
	ctx := context.Background()

	// Obtener todos los agentes
	agents, err := repo.GetAllAgents(ctx)
	if err != nil {
		log.Printf("[Scheduler] Error obteniendo agentes: %v", err)
		return
	}

	log.Printf("[Scheduler] Procesando %d agentes...", len(agents))

	for _, agent := range agents {
		// Verificar si el agente está en cooldown
		if time.Now().Before(agent.CooldownUntil) {
			log.Printf("[Scheduler] Agente %s en cooldown hasta %v, saltando...", agent.ID, agent.CooldownUntil)
			continue
		}

		// Ejecutar tick para este agente
		log.Printf("[Scheduler] Ejecutando tick para agente %s...", agent.ID)
		if err := engine.RunTick(ctx, agent.ID); err != nil {
			log.Printf("[Scheduler] Error en tick del agente %s: %v", agent.ID, err)
		}
	}
}
