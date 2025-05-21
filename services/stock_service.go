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
	// Validate date range
	if startDate.After(endDate) {
		return nil, errors.NewModelValidationError(
			"StockService",
			"date_range",
			"start date cannot be after end  date",
		)
	}

	// Get stock prices from the repository
	stocks, err := s.stockRepo.RetrieveStocksFromDatabase(ctx, symbol, startDate, endDate)
	if err != nil {
		return nil, errors.NewServiceError("Retrieving stock history", err)
	}

	return stocks, nil
}
