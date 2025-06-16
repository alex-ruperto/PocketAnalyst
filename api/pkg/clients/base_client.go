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

// Common logic for making HTTP requests to an API, with all the error handling.
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

// MakeArrayRequest handles HTTP requests that return JSON arrays instead of objects.
// Some APIs (like FMP) return arrays directly rather than wrapping them in objects.
func (bc *BaseClient) MakeArrayRequest(url string) ([]map[string]any, error) {
	// Make HTTP request with proper error wrapping
	resp, err := bc.Client.Get(url)
	if err != nil {
		return nil, client_errors.NewHTTPRequestError(url, err)
	}
	defer resp.Body.Close() // Ensure response body is closed to prevent resource leaks

	// Check for successful HTTP response status
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, client_errors.NewHTTPStatusError(url, resp.StatusCode, string(body))
	}

	// Read the entire response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, client_errors.NewResponseReadError(err)
	}

	// Parse JSON response as an array
	var response []map[string]any
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, client_errors.NewResponseParseError(err)
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

// CheckArrayAPIError examines array responses for provider-specific error messages.
// Some APIs return errors in the first element of an array response.
func (bc *BaseClient) CheckArrayAPIError(response []map[string]any, errorKeys ...string) error {
	// Check for empty response first
	if len(response) == 0 {
		return client_errors.NewAPIError("no data returned - possibly invalid symbol or API error")
	}

	// Default error keys to check in the first array element
	defaultKeys := []string{"error", "message", "error_message", "Error Message"}

	// Combine provided keys with defaults
	allKeys := append(errorKeys, defaultKeys...)

	// Check the first element for error messages
	firstElement := response[0]
	for _, key := range allKeys {
		if errorMsg, exists := firstElement[key]; exists {
			if msg, ok := errorMsg.(string); ok && msg != "" {
				return client_errors.NewAPIError(msg)
			}
		}
	}

	return nil
}
