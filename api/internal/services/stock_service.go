package services

import (
	"context"
	"fmt"
	"pocketanalyst/internal/models"
	"pocketanalyst/internal/repositories"
	"pocketanalyst/pkg/clients"
	"pocketanalyst/pkg/errors"
	"strings"
	"time"
)

// StockService handles business logic related to stock operations
type StockService struct {
	stockRepo *repositories.StockRepository
	client    clients.StockDataClient
}

func NewStockService(
	stockRepo *repositories.StockRepository,
	client clients.StockDataClient,
) *StockService {
	return &StockService{
		stockRepo: stockRepo,
		client:    client,
	}
}

func (s *StockService) SynchronizeStockData(ctx context.Context, symbol string) (int, error) {
	// Fetch stock data from chosen API
	stocks, err := s.client.FetchDaily(symbol)
	if err != nil {
		return 0, errors.NewServiceError("Fetching stock data", err)
	}

	// If no data was returned, return early
	if len(stocks) == 0 {
		return 0, fmt.Errorf("No stock data found for symbol %s", symbol)
	}

	// Store the fetched data in the database
	storedCount, err := s.stockRepo.SaveStocksToDatabase(ctx, stocks)
	if err != nil {
		return 0, errors.NewServiceError("Storing stock data", err)
	}

	return storedCount, nil
}

func (s *StockService) GetStockHistory(
	ctx context.Context,
	symbol string,
	startDate, endDate time.Time,
) ([]*models.Stock, error) {
	// Validate date function parameters before anything
	if err := s.validateInput(symbol, startDate, endDate); err != nil {
		return nil, err
	}

	// Get stock prices from the database using the repository
	stocks, err := s.stockRepo.RetrieveStocksFromDatabase(ctx, symbol, startDate, endDate)
	if err != nil {
		return nil, errors.NewServiceError("Retrieving stock history", err)
	}

	// If no data is returned, treat as "symbol not found."
	if len(stocks) == 0 {
		return nil, errors.NewNotFoundError("Symbol", symbol)
	}
	return stocks, nil
}

func (s *StockService) validateInput(symbol string, startDate, endDate time.Time) error {
	if strings.TrimSpace(symbol) == "" {
		return errors.NewModelValidationError(
			"StockService",
			"symbol",
			"symbol cannot be empty",
		)
	}

	if startDate.After(endDate) {
		return errors.NewModelValidationError(
			"StockService",
			"date_range",
			"start date cannot be after end date",
		)
	}

	if startDate.IsZero() || endDate.IsZero() {
		return errors.NewModelValidationError(
			"StockService",
			"date_range",
			"dates cannot be empty",
		)
	}

	return nil
}
