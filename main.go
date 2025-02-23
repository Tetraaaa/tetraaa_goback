package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"tetraaa/goback/utils"
	"time"

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

func httpError(w http.ResponseWriter, errorMessage string) {
	type Response struct {
		Error string `json:"error"`
	}
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	response := Response{Error: errorMessage}
	jsonStr, _ := json.Marshal(response)
	w.Write(jsonStr)
}

func startHttpServer() {
	var startTime = time.Now()

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

	http.HandleFunc("/status", func(w http.ResponseWriter, req *http.Request) {

		temp := utils.GetCPUTemp()
		avgs := utils.GetCPUSAverages()
		memTotal, memFree := utils.GetMemoryTotalAndFree()
		peribotStatus, _ := utils.GetPeribotStatus()

		type Response struct {
			CpuTemp       float64               `json:"cpuTemp"`
			Uptime        time.Duration         `json:"uptime"`
			CPUSAverages  []int64               `json:"cpus"`
			MemTotal      uint64                `json:"memTotal"`
			MemFree       uint64                `json:"memFree"`
			PeribotStatus utils.PeribotResponse `json:"peribotStatus"`
		}
		response := Response{CpuTemp: temp, Uptime: time.Duration(time.Since(startTime).Seconds()), CPUSAverages: avgs, MemTotal: memTotal, MemFree: memFree, PeribotStatus: peribotStatus}
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
	// databaseConnection := connectToDatabase()
	// var result string
	// err := databaseConnection.QueryRow(context.Background(), "select 'Hello, world!'").Scan(&result)
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
	// 	os.Exit(1)
	// }
	startHttpServer()
	// defer databaseConnection.Close()
}
