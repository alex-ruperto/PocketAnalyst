package clients

import (
	"pocketanalyst/internal/models"
)

// This interface will allow us to swap between data sources as needed.
type StockDataClient interface {
	FetchDaily(symbol string) ([]*models.Stock, error)
	GetProviderName() string
}

// Holds configuration parameters for creating clients.
type ClientConfig struct {
	BaseURL string `json:"base_url"`
	APIKey  string `json:"api_key"`
}
