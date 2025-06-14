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

type FMPClient struct {
	BaseURL string       // API base URL
	APIKey  string       // API key for authentication
	Client  *http.Client // HTTP client with timeouts
}

func NewFMPClient(baseURL, apiKey string) *FMPClient {
	return &FMPClient{
		BaseURL: baseURL,
		APIKey:  apiKey,
		Client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (fmpc *FMPClient) FetchDailyPricesFromFMP(symbol string) ([]*models.Stock, error) {
	url := fmt.Sprintf("%s/full?symbol=%s&apikey=%s", fmpc.BaseURL, symbol, fmpc.APIKey)

	// Make the HTTP Request. Return HTTPRequestErorr if it fails.
	resp, err := fmpc.Client.Get(url)
	if err != nil {
		return nil, client_errors.NewHTTPRequestError(url, err)
	}

	// Ensure response body is closed to prevent a resource leak.
	defer resp.Body.Close()

	// Check to see that the HTTP response was successful. HTTPStatusError if it was not a successful response.
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, client_errors.NewHTTPStatusError(url, resp.StatusCode, string(body))
	}

	// Read response body. Return a ResponseReadError if it fails.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, client_errors.NewResponseReadError(err)
	}

	// Parse JSON Response and store it into a map.
	var response map[string]any

	// Use unmarshal to read the body. Pass response by reference here to modify the original.
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, client_errors.NewResponseParseError(err)
	}

}
