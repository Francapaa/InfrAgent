package main

import (
	"log"
	"server/internal/database"
	
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error cargando .env")
	}

	database.ConnectMongo()
		
}
