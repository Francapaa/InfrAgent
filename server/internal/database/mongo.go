package database

import (
	"context"
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

var DB *pgxpool.Pool
var SQLDB *sql.DB

func ConnectDatabase() *pgxpool.Pool {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error cargando la base de datos")
	}

	connStr := os.Getenv("DATABASE_URL")

	if connStr == "" {
		log.Fatal("No existe esa variable de entorno")
	}

	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		panic(err)
	}
	config.MaxConns = 10
	config.MinConns = 2
	config.MaxConnLifetime = time.Hour

	db, err := pgxpool.NewWithConfig(context.Background(), config)

	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(context.Background()); err != nil {
		log.Fatal(err)
	}

	sqlDB := stdlib.OpenDB(*db.Config().ConnConfig)
	if sqlDB == nil {
		log.Fatal("Error creando sql.DB")
	}
	SQLDB = sqlDB

	log.Println("Connected with PostgreSQL (Neon)")
	DB = db
	return db
}

func GetSQLDB() *sql.DB {
	return SQLDB
}
