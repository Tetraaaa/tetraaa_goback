package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"tetraaa/goback/finance"
	"tetraaa/goback/utils"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func loadEnvVars() {
	hostname, _ := os.Hostname()
	env := "dev"
	if hostname == "raspberrypi" {
		env = "prod"
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
	w.WriteHeader(http.StatusInternalServerError)
	response := Response{Error: errorMessage}
	jsonStr, _ := json.Marshal(response)
	w.Write(jsonStr)
}

func authError(w http.ResponseWriter) {
	type Response struct {
		Error string `json:"error"`
	}
	w.WriteHeader(http.StatusUnauthorized)
	response := Response{Error: "Invalid credentials"}
	jsonStr, _ := json.Marshal(response)
	w.Write(jsonStr)
}

func checkForAuth(req *http.Request, w http.ResponseWriter) bool {
	validApiKey := os.Getenv("API_KEY")
	if validApiKey == "" {
		return true
	}

	authHeader := req.Header.Get("Authorization")
	if authHeader == "" {
		authError(w)
		return false
	}

	if authHeader == fmt.Sprintf("Bearer %s", validApiKey) {
		return true
	}

	authError(w)
	return false
}

func route(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	http.HandleFunc(pattern, func(w http.ResponseWriter, req *http.Request) {
		w.Header().Add("Access-Control-Allow-Headers", " Origin, Content-Type, Authorization")
		w.Header().Add("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Content-Type", "application/json")
		handler(w, req)
	})
}

func startHttpServer() {
	var startTime = time.Now()

	route("/ping", func(w http.ResponseWriter, req *http.Request) {
		type Response struct {
			Message string `json:"message"`
		}
		response := Response{Message: "Pong !"}
		jsonStr, _ := json.Marshal(response)
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Content-Type", "application/json")
		w.Write(jsonStr)
	})

	route("/portfolio", func(w http.ResponseWriter, req *http.Request) {
		if !checkForAuth(req, w) {
			return
		}
		portfolio, err := finance.GetPortfolio()
		if err != nil {
			httpError(w, "Unable to retrieve portfolio. Make sure a data/portfolio.json file exists and retry.")
			return
		}

		jsonStr, _ := json.Marshal(portfolio)
		w.Write(jsonStr)
	})

	route("/portfolio-history", func(w http.ResponseWriter, req *http.Request) {
		if !checkForAuth(req, w) {
			return
		}
		type Response struct {
			PortfolioHistory finance.PortfolioHistory `json:"portfolio_history"`
		}

		portfolio_history, err := os.ReadFile("data/portfolio_history.json")
		if err != nil {
			httpError(w, "Unable to retrieve portfolio history. Make sure a data/portfolio_history.json file exists and retry.")
			return
		}

		response := Response{}
		json.Unmarshal(portfolio_history, &response)
		jsonStr, _ := json.Marshal(response)
		w.Write(jsonStr)
	})

	route("/stock", func(w http.ResponseWriter, req *http.Request) {
		if !checkForAuth(req, w) {
			return
		}
		if !req.URL.Query().Has("ticker") {
			httpError(w, "Missing ticker query param")
			return
		}

		ticker := req.URL.Query().Get("ticker")

		stockData, err := finance.GetStockInformationFromTicker(ticker)
		if err != nil {
			httpError(w, "Unable to get stock data for ticker")
			return
		}

		jsonStr, _ := json.Marshal(stockData)
		w.Write(jsonStr)

	})

	route("/status", func(w http.ResponseWriter, req *http.Request) {

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
		w.Write(jsonStr)
	})

	fmt.Println("Démarrage du serveur sur le port 8080.")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Impossible de démarrer le serveur http : ", err)
	}
}

func main() {
	fmt.Println("Démarrage...")
	fmt.Println("Chargement des variables d'environnement...")
	loadEnvVars()
	fmt.Println("Enregistrement des cron jobs...")
	utils.RegisterCronJobs()

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
