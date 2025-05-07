package repositories

import (
	"GoSimpleREST/models"
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

	// If commit doesn't occur, rollback will occur.
	// this will be ignored if tx.Commit is called
	defer tx.Rollback()

	// Process each stock
	for _, stock := range stocks {
		// Validate the stock data using model's validation
		if err := stock.Validate(); err != nil {
			return 0, err
		}

		// Check if company exists in the companies table
		var companyID int
		err := tx.QueryRowContext(
			ctx,
			`SELECT company_id FROM companies WHERE symbol = $1`,
			stock.Symbol,
		).Scan(&companyID)

		// If the company doesn't exist, create it
		if err == sql.ErrNoRows {
			err = tx.QueryRowContext(
				ctx,
				`INSERT INTO companies (symbol, name, is_active, created_at)
				VALUES ($1, $1, true, NOW())
				RETURNING company_id`,
				stock.Symbol,
			).Scan(&companyID)

			if err != nil {
				return 0, fmt.Errorf("Failed to create company for symbol %s: %w",
					stock.Symbol, err)

			}
		} else if err != nil {
			return 0, fmt.Errorf("Failed to check if company exists for symbol %s: %w",
				stock.Symbol, err)
		}

		// Set the company ID for the stock
		stock.CompanyID = companyID
	}

	// Use a prepared statement for efficient batch insertion
	// Helps prevent SQL injection
	// No need to include price_id PK here since it auto-increments (SERIAL type auto-incrementing PK)
	// More efficient than	constructing the query for each stock
	// Defer stmt.Close() to ensure the statement is properly closed when we're done to prevent resource leaks.
	stmt, err := tx.PrepareContext(
		ctx,
		`INSERT INTO stock_prices
		(company_id, symbol, date, open_price, high_price, low_price,
		close_price, adjusted_close, volume, dividend_amount, 
		split_coefficient, source_id, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, NOW())
		ON CONFLICT (company_id, date, source_id)
		DO UPDATE SET
		open_price = EXCLUDED.open_price,
		high_price = EXCLUDED.high_price,
		low_price = EXCLUDED.low_price,
		close_price = EXCLUDED.close_price,
		adjusted_close = EXCLUDED.adjusted_close,
		volume = EXCLUDED.volume,
		dividend_amount = EXCLUDED.dividend_amount,
		split_coefficient = EXCLUDED.split_coefficient,
		created_at = EXCLUDED.created_at
		RETURNING price_id
		`,
	)
	if err != nil {
		return 0, fmt.Errorf("Failed to prepare stock price insert statement: %w", err)
	}

	defer stmt.Close()

	for _, stock := range stocks {
		// Source ID for Alpha Vantage (assuming it's 1)
		// TODO: make this a more flexible implementation to dynamically select the correct source ID.
		const sourceID = 1

		// Execute the prepared statement with values for this stock
		// RETURNING price_id gives us back the auto-generated PK
		var priceID int
		err = stmt.QueryRowContext(
			ctx,
			stock.CompanyID,
			stock.Symbol,
			stock.Date,
			stock.OpenPrice,
			stock.HighPrice,
			stock.LowPrice,
			stock.ClosePrice,
			stock.AdjustedClosePrice,
			stock.Volume,
			stock.DividendAmount,
			stock.SplitCoefficient,
			sourceID,
		).Scan(&priceID)

		// Handle the error
		if err != nil {
			return 0, fmt.Errorf("failed to insert stock price for %s on %s: %w",
				stock.Symbol, stock.Date.Format("2006-01-02"), err)
		}

		// Update the stock object with the Generated ID
		// This keeps our in-memory data in sync with the DB.
		stock.PriceID = priceID
	}

	// Commit the transaction
	// This makes all our changes permanent in the DB.
	// If this fails, the deferred rollback will undo everything.
	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("Failed to commit transaction: %w", err)
	}

	return len(stocks), nil
}
