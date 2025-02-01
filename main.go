package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func loadEnvVars() {
	env := os.Getenv("ENV")
	if env == "" {
		env = "dev"
	}
	err := godotenv.Load(".env." + env)
	if err != nil {
		log.Fatal("Impossible de charger le .env : ", err)
	}
}

func connectToDatabase() *pgxpool.Pool {
	conn, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Impossible de se connecter à la base de données : ", err)
	}
	return conn
}

func main() {
	fmt.Println("Démarrage...")
	loadEnvVars()
	databaseConnection := connectToDatabase()
	var result string
	err := databaseConnection.QueryRow(context.Background(), "select 'Hello, world!'").Scan(&result)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}
	defer databaseConnection.Close()
}
