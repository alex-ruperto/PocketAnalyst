package clients

import (
	"encoding/json"
	"io"
	"net/http"
	"pocketanalyst/pkg/errors/client_errors"
	"time"
)

// Struct for common HTTP Client functionality. Promotes code reuse.
type BaseClient struct {
	BaseURL string
	APIKey  string
	Client  *http.Client
}

func NewBaseClient(baseURL, apiKey string) *BaseClient {
	return &BaseClient{
		BaseURL: baseURL,
		APIKey:  apiKey,
		Client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (bc *BaseClient) MakeRequest(url string) (map[string]any, error) {
	// Make HTTP request. Return a HTTPRequestError if it fails.
	resp, err := bc.Client.Get(url)
	if err != nil {
		return nil, client_errors.NewHTTPRequestError(url, err)
	}
	defer resp.Body.Close() // Ensure the response body is closed to prevent resource leaks.

	// Check for successful HTTP response. Return a HTTPStatusError if it is not a successful response.
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, client_errors.NewHTTPStatusError(url, resp.StatusCode, string(body))
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, client_errors.NewResponseReadError(err)
	}

	// Parse the JSON response
	var response map[string]any
	// Unmarshal only needs to read body. We need to modify the original response so we pass response by reference here.
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, client_errors.NewResponseParseError(err)
	}

	// Check for API error messages
	if errorMsg, exists := response["Error Message"]; exists {
		return nil, client_errors.NewAPIError(errorMsg.(string))
	}

	return response, nil
}

// NOTE: The 'errorsKeys ...string' means accept zero or more string elements. Variadic parameter.

// Function to check specific API errors from a client.
func (bc *BaseClient) CheckAPIError(response map[string]any, errorKeys ...string) error {
	// Default error keys to check - covers most common API error patterns
	defaultKeys := []string{"Error Message", "error", "message", "error_message"}

	// NOTE: The 'defaultKeys...' means append each element of defaultKeys to errorKeys (instead of appending the entire slice to one element).
	allKeys := append(errorKeys, defaultKeys...)

	// Check each possible error key
	for _, key := range allKeys {
		if errorMsg, exists := response[key]; exists {
			if msg, ok := errorMsg.(string); ok && msg != "" {
				return client_errors.NewAPIError(msg)
			}
		}
	}

	return nil
}
