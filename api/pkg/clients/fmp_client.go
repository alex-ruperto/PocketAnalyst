package clients

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"pocketanalyst/internal/models"
	"pocketanalyst/pkg/errors/client_errors"
	"sort"
	"strconv"
	"time"
)

type FMPCient struct {
	BaseURL string       // API base URL
	APIKey  string       // API key for authentication
	Client  *http.Client // HTTP client with timeouts
}

func FMPClient(baseURL, apiKey string) *AlphaVantageClient {
	return &AlphaVantageClient{
		BaseURL: baseURL,
		APIKey:  apiKey,
		Client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}
