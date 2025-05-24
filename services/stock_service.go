package services

import (
	"PocketAnalyst/clients"
	"PocketAnalyst/errors"
	"PocketAnalyst/models"
	"PocketAnalyst/repositories"
	"context"
	"fmt"
	"time"
)

// StockService handles business logic related to stock operations
type StockService struct {
	stockRepo   *repositories.StockRepository
	alphaClient *clients.AlphaVantageClient
}

func NewStockService(
	stockRepo *repositories.StockRepository,
	alphaClient *clients.AlphaVantageClient,
) *StockService {
	return &StockService{
		stockRepo:   stockRepo,
		alphaClient: alphaClient,
	}
}

func (s *StockService) SynchronizeStockData(ctx context.Context, symbol string) (int, error) {
	// Fetch stock data from Alpha Vantage API
	stocks, err := s.alphaClient.FetchDailyPricesFromAPI(symbol)
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

	return stocks, nil
}

func (s *StockService) validateInput(symbol string, startDate, endDate time.Time) error {
	if symbol == "" {
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
