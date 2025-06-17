package main

import (
	"log"
	"os"
	"pocketanalyst/internal/app"
	"strconv"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver
)

func main() {
	// Load configuration from environment variables
	config := loadConfig()

	// Create and initalize the application
	app, err := app.NewApp(config)
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	// Ensure a graceful shutdown
	defer func() {
		if err := app.Close(); err != nil {
			log.Printf("Error during application shutdown: %v", err)
		}
	}()

	// Start the server
	if err := app.Start(); err != nil {
		log.Fatalf("Server failed to start %v", err)
	}
	log.Printf("Server started successfully.")
}

func loadConfig() *app.Config {
	dbURL := getEnvWithDefault("DATABASE_URL", "postgres://localhost/pocketanalyst?sslmode=disable")

	return &app.Config{
		DatabaseURL:           dbURL,
		FMPAPIKey:             getEnvWithDefault("FMP_API_KEY", ""),
		FMPBaseURL:            getEnvWithDefault("FMP_BASE_URL", "https://financialmodelingprep.com"),
		Port:                  getEnvWithDefault("PORT", "8080"),
		ReadTimeout:           time.Duration(getEnvAsInt("READ_TIMEOUT_SECONDS", 30)) * time.Second,
		WriteTimeout:          time.Duration(getEnvAsInt("WRITE_TIMEOUT_SECONDS", 30)) * time.Second,
		MaxIdleConnections:    getEnvAsInt("MAX_IDLE_CONNECTIONS", 10),
		MaxOpenConnections:    getEnvAsInt("MAX_OPEN_CONNECTIONS", 100),
		ConnectionMaxLifetime: time.Duration(getEnvAsInt("CONNECTION_MAX_LIFETIME_MINUTES", 60)) * time.Minute,
	}
}

// getEnvWithDefault returns environment variables as strings or default if not set/invalid
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
