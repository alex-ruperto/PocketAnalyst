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

func (s *StockService) FetchAndStoreStockData(ctx context.Context, symbol string) (int, error) {
	// Fetch stock data from Alpha Vantage
	stocks, err := s.alphaClient.FetchDaily(symbol)
	if err != nil {
		return 0, errors.NewServiceError("fetching stock data", err)
	}

	// If no data was returned, return early
	if len(stocks) == 0 {
		return 0, fmt.Errorf("no stock data found for symbol %s", symbol)
	}

	// Store the fetched data in the database
	storedCount, err := s.stockRepo.StoreStockPrices(ctx, stocks)
	if err != nil {
		return 0, errors.NewServiceError("storing stock data", err)
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
	stocks, err := s.stockRepo.GetStockPrices(ctx, symbol, startDate, endDate)
	if err != nil {
		return nil, errors.NewServiceError("retrieving stock history", err)
	}

	return stocks, nil
}
