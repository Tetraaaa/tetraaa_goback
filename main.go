package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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

func startHttpServer() {
	http.HandleFunc("/ping", func(w http.ResponseWriter, req *http.Request) {
		type Response struct {
			Message string `json:"message"`
		}
		response := Response{Message: "Pong !"}
		jsonStr, _ := json.Marshal(response)
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Content-Type", "application/json")
		w.Write(jsonStr)
	})
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Impossible de démarrer le serveur http : ", err)
	}
	fmt.Println("Serveur démarré sur le port 8080...")
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
	startHttpServer()
	defer databaseConnection.Close()
}
