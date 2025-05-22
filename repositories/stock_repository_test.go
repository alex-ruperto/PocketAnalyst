package repositories

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver
)

func TestStockRepository_CheckIfSymbolExists(t *testing.T) {
	// Set up DB the connection string
	dbConnectionString := os.Getenv("DATABASE_URL")
	if dbConnectionString == "" {
		t.Skip("DATABASE_URL not set. Skipping tests.")
	}

	// Connect to DB
	db, err := sql.Open("postgres", dbConnectionString)
	if err != nil {
		t.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()

	// Verify connection works
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		t.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("Connected to database successfully")

	stockRepo := NewStockRepository(db)

	t.Run("Check existing symbol", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		exists, err := stockRepo.CheckIfSymbolExists(ctx, "AAPL")
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		t.Logf("Symbol AAPL exists: %v", exists)
	})
}
