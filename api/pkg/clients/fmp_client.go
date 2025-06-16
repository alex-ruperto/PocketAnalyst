package clients

import (
	"fmt"
	"pocketanalyst/internal/models"
	"strconv"
	"time"
)

type FMPClient struct {
	*BaseClient // Embed BaseClient functionality into FMPClient
}

func NewFMPClient(baseURL, apiKey string) *FMPClient {
	return &FMPClient{
		BaseClient: NewBaseClient(baseURL, apiKey),
	}
}

func (fmpc *FMPClient) GetProviderName() string {
	return "FMP"
}

func (fmpc *FMPClient) FetchDaily(symbol string) ([]*models.Stock, error) {
	url := fmt.Sprintf("%s/stable/historical-price-eod/full?symbol=%s&apikey=%s",
		fmpc.BaseURL, symbol, fmpc.APIKey)

	// Use the shared HTTP Request logic from BaseClient
	dailyData, err := fmpc.MakeArrayRequest(url)
	if err != nil {
		return nil, err
	}

	// Check for FMP-specific error messages.
	if err := fmpc.CheckArrayAPIError(dailyData); err != nil {
		return nil, err
	}

	// Convert to Stock models
	stocks := make([]*models.Stock, 0, 5)
	for i, dayData := range dailyData {
		if i >= 5 { // Only get first 5 records
			break
		}

		// Parse date
		dateStr, ok := dayData["date"].(string)
		if !ok {
			continue // Skip if no valid date
		}

		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			continue // Skip if date parsing fails
		}

		// Create Stock model
		stock := &models.Stock{
			Symbol:           symbol,
			Date:             date,
			OpenPrice:        getFloat(dayData, "open"),
			HighPrice:        getFloat(dayData, "high"),
			LowPrice:         getFloat(dayData, "low"),
			ClosePrice:       getFloat(dayData, "close"),
			AdjustedClose:    getFloat(dayData, "close"), // FMP doesn't provide adjusted_close in this endpoint
			Volume:           getFloat(dayData, "volume"),
			DividendAmount:   0, // Not available in this endpoint
			SplitCoefficient: 1, // Not available in this endpoint
			DataSource:       fmpc.GetProviderName(),
			LastUpdated:      time.Now(),
		}

		stocks = append(stocks, stock)
	}

	return stocks, nil
}

// getFloat handles handles FMP's numeric format
func getFloat(data map[string]any, key string) float64 {
	if val, ok := data[key]; ok {
		switch v := val.(type) {
		case float64:
			return v
		case int:
			return float64(v)
		case int64:
			return float64(v)
		case string:
			// Fallback to string parsing if FMP sometimes returns strings
			if f, err := strconv.ParseFloat(v, 64); err == nil {
				return f
			}
		}
	}
	return 0.0
}
