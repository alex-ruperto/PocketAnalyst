package repositories

import (
	"GoSimpleRest/models"
	"context"
	"database/sql"
	"fmt"
	"time"
)

// StockRepository handles database operations for stocks
type StockRepository struct {
	db *sql.DB
}

// NewStockRepository creates a new stock repository
func NewStockRepository(db *sql.DB) *StockRepository {
	return &StockRepository{db: db}
}

// StoreStockPrices stores multiple stock price records
func (sr *StockRepository) StoreStockPrices(ctx context.Context, stocks []*models.Stock) (int, error) {
	// Begin transaction
	tx, err := sr.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback() // this will be ignored if tx.Commit is called

	// Process each stock
	for _, stock := range stocks {
		// Validate the stock data using model's validation
		if err := stock.Validate(); err != nil {
			return 0, err
		}

		// Check if company exists
		var companyID int
	}
}
