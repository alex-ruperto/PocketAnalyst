package main

import (
	"PocketAnalyst/api"
	"PocketAnalyst/clients"
	"PocketAnalyst/controllers"
	"PocketAnalyst/repositories"
	"PocketAnalyst/services"
	"database/sql"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver
)

func main() {
	// Get configuration from env variables
	dbConnectionString := os.Getenv("DATABASE_URL")
	alphaVantageAPIKey := os.Getenv("ALPHA_VANTAGE_API_KEY")
	alphaVantageBaseURL := os.Getenv("ALPHA_VANTAGE_BASE_URL")
	port := "8080"

	// Connect to DB
	db, err := sql.Open("postgres", dbConnectionString)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	defer db.Close()

	// Verify db connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("Connected to database successfully")

	// Create instances of all components
	stockRepo := repositories.NewStockRepository(db)
	alphaClient := clients.NewAlphaVantageClient(alphaVantageBaseURL, alphaVantageAPIKey)
	stockService := services.NewStockService(stockRepo, alphaClient)
	stockController := controllers.NewStockController(stockService)

	// Set up HTTP Routes
	http.HandleFunc("/api/stocks/fetch", stockController.HandleStockFetchRequest)
	http.HandleFunc("/api/stocks/get", stockController.HandleStockHistoryRequest)

	// Start HTTP Server
	log.Printf("Starting server on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func loadConfig() *api.Config {
	return &api.Config{
		DatabaseURL:           getEnvWithDefault("DATABASE_URL", "postgres://localhost/pocketanalyst?sslmode=disable"),
		AlphaVantageAPIKey:    getEnvWithDefault("ALPHA_VANTAGE_API_KEY", ""),
		AlphaVantageBaseURL:   getEnvWithDefault("ALPHA_VANTAGE_BASE_URL", "https://www.alphavantage.co/query"),
		Port:                  getEnvWithDefault("PORT", "8080"),
		ReadTimeout:           time.Duration(getEnvAsInt("READ_TIMEOUT_SECONDS", 30)) * time.Second,
		WriteTimeout:          time.Duration(getEnvAsInt("WRITE_TIMEOUT_SECONDS", 30)) * time.Second,
		MaxIdleConnections:    getEnvAsInt("MAX_IDLE_CONNECTIONS", 10),
		MaxOpenConnections:    getEnvAsInt("MAX_OPEN_CONNECTIONS", 100),
		ConnectionMaxLifetime: time.Duration(getEnvAsInt("CONNECTION_MAX_LIFETIME_MINUTES", 60)) * time.Minute,
	}
}

func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt returns environment variables as integer or default if not set/invalid
func getEnvAsInt(key string, defaultValue int) int {
	if valueStr := os.Getenv(key); valueStr != "" {
		if value, err := strconv.Atoi(valueStr); err == nil {
			return value
		}
		log.Printf("Warning: Invalid integer value for %s: %s, using default: %d", key, valueStr, defaultValue)
	}
	return defaultValue
}
