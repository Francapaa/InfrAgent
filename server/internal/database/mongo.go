package database

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client
var MongoDB *mongo.Database

func ConnectMongo() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	uri := os.Getenv("MONGO_URI")
	dbName := os.Getenv("MONGO_DB")

	log.Printf("🔍 URI: %s", uri)
	log.Printf("🔍 DB Name: %s", dbName)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal("Error conectando a MongoDB:", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("MongoDB no responde:", err)
	}

	MongoClient = client
	MongoDB = client.Database(dbName)

	log.Println("✅ MongoDB conectado correctamente")
}

func GetDB() *mongo.Database {
	return MongoDB
}

// para cerrar la conexion a la base de datos
func DisconnectMongo() {

	if MongoClient != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := MongoClient.Disconnect(ctx); err != nil {
			log.Println("Error desconectando de MongoDB:", err)
		}
	}

}
