package main

import (
	"PocketAnalyst/clients"
	"PocketAnalyst/controllers"
	"PocketAnalyst/repositories"
	"PocketAnalyst/services"
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
	http.HandleFunc("/api/stocks/get"), stockController.HandleStockHistoryRequest)

	// Start HTTP Server
	log.Printf("Starting server on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
