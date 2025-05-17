package main

import (
	"PocketAnalyst/clients"
	"PocketAnalyst/controllers"
	"PocketAnalyst/repositories"
	"PocketAnalyst/service"
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq" // PostgreSQL driver
)

func main() {
	// Get configuration from env variables
	dbConnectionString := os.Getenv("DATABASE_URL")
	alphaVantageAPIKey := os.Getenv("ALPHA_VANTAGE_API_KEY")
	alphaVantageBaseURL := os.Getenv("ALPHA_VANTAGE_BASE_URL")
	port := 8080
}
