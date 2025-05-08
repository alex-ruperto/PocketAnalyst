package clients

import (
	"GoSimpleREST/models"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type AlphaVantageClient struct {
	BaseURL string       // API base URL
	APIKey  string       // API key for authentication
	Client  *http.Client // HTTP client with timeouts
}

func NewAlphaVantageClient(baseURL, apiKey string) *AlphaVantageClient {
	return &AlphaVantageClient{
		BaseURL: baseURL,
		APIKey:  apiKey,
		Client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}
