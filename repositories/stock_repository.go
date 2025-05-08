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

// StoreStockPrices stores multiple stock price records. Using a transaction, all operations will
// either succeed or fail together.
func (sr *StockRepository) StoreStockPrices(ctx context.Context, stocks []*models.Stock) (int, error) {
	// Begin transaction
	// Transaction will ensures all stock prices are inserted/updated or none are.
	// This preserves data consistency in case of errors.
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
				`
				INSERT INTO companies (symbol, name, is_active, created_at)
				VALUES ($1, $1, true, NOW())
				RETURNING company_id
				`,
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

	// Use a prepared statement for efficient batch insertion as it can be reused.
	// Helps prevent SQL injection
	// No need to include price_id PK here since it auto-increments (SERIAL type auto-incrementing PK)
	// More efficient than	constructing the query for each stock
	// Defer stmt.Close() to ensure the statement is properly closed when we're done to prevent resource leaks
	// ON CONFLICT statement will prevent duplicate entries and EXCLUDED uses the value from row tried to insert
	stmt, err := tx.PrepareContext(
		ctx,
		`
		INSERT INTO stock_prices
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

	// Close prepared statement when we are done with it.
	defer stmt.Close()

	// Insert each stock price.
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
			stock.AdjustedClose,
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

// Retrieves stock prices for a symbol within a date range. The date range helps limit the data returned to what
// is actually needed.
func (sr *StockRepository) GetStockPrices(
	ctx context.Context,
	symbol string,
	startDate, endDate time.Time,
) ([]*models.Stock, error) {
	// SQL query with JOIN to get data source name
	query := `
		SELECT sp.price_id, sp.company_id, sp.symbol, sp.date, 
		       sp.open_price, sp.high_price, sp.low_price, sp.close_price, 
		       sp.adjusted_close, sp.volume, sp.dividend_amount, 
		       sp.split_coefficient, ds.source_name, sp.created_at
		FROM stock_prices sp
		JOIN data_sources ds ON sp.source_id = ds.source_id
		WHERE sp.symbol = $1 AND sp.date BETWEEN $2 AND $3
		ORDER BY sp.date DESC
	`

	// Execute the query with parameters
	// QueryContext ensures the query can be cancelled if the context is cancelled
	rows, err := sr.db.QueryContext(ctx, query, symbol, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to query stock prices: %w", err)
	}

	// Always close rows when done to prevent resource leaks
	defer rows.Close()

	// Process each row returned by the query
	var stocks []*models.Stock
	for rows.Next() {
		var s models.Stock
		var createdAt time.Time

		// Scan row data into struct fields
		// The order here has to match that of the SELECT statement.
		err := rows.Scan(
			&s.PriceID,
			&s.CompanyID,
			&s.Symbol,
			&s.Date,
			&s.OpenPrice,
			&s.HighPrice,
			&s.LowPrice,
			&s.ClosePrice,
			&s.AdjustedClose,
			&s.Volume,
			&s.DividendAmount,
			&s.SplitCoefficient,
			&s.DataSource,
			&createdAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan stock price row: %w", err)
		}

		s.CreatedAt = createdAt

		// Return the stocks as a slice
		stocks = append(stocks, &s)
	}

	// Check for errors from iterating rows
	// This catches any errors that occurred during iteration
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating stock price rows: %w", err)
	}

	return stocks, nil
}
