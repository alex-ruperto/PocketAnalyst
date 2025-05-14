package clients

import (
	"PocketAnalyst/errors/client_errors"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestAlphaVantageClient_FetchDaily(t *testing.T) {
	if os.Getenv("RUN_INTEGRATION_TESTS") != "true" {
		t.Skip("Skipping integration test; set RUN_INTEGRATION_TESTS=true to run.")
	}

	// Get the API Key from environment variable
	apiKey := os.Getenv("ALPHA_VANTAGE_API_KEY")
	if apiKey == "" {
		t.Fatal("ALPHA_VANTAGE_API_KEY environment variable is not set")
	}

	// Create a real Alpha Vantage client
	client := NewAlphaVantageClient("https://www.alphavantage.co/query", apiKey)

	// Make the actual API call
	stocks, err := client.FetchDaily("AAPL")
	if err != nil {
		t.Fatalf("Failed to fetch data from Alpha Vantage: %v", err)
	}

	// Basic verification of the response
	if len(stocks) == 0 {
		t.Error("No stocks data was returned")
	}

	t.Log("Displaying first 5 stock data points")
	for i, stock := range stocks {
		if i >= 5 {
			break
		}
		t.Logf("Stock #%d: %+v", i+1, stock)
	}

}

func TestAlphaVantageClient_FetchDaily_Errors(t *testing.T) {
	// Test case 1: HTTP Request failure
	t.Run("HTTP Request failure", func(t *testing.T) {
		client := NewAlphaVantageClient("http://non-existent-url.example", "dummy-api-key")
		stocks, err := client.FetchDaily("IBM")

		if err == nil {
			t.Error("Expected an error for a non-existent URL, but got nil.")
		}
		if stocks != nil {
			t.Error("Expected nil stocks, but got some data.")
		}

		expectedError := &client_errors.HTTPRequestError{}
		expectedCode := expectedError.ErrorCode()

		// Check error type
		if !errors.As(err, &expectedError) {
			t.Errorf("Expected error of type *client_errors.HTTPRequestError, but got %T", err)
		}

		// Check interface implementation and error code
		clientErr, ok := err.(client_errors.ClientError)
		if !ok {
			t.Errorf("Error does not implement ClientError interface")
		} else if clientErr.ErrorCode() != expectedCode {
			t.Errorf("Expected error code %s, but got %s", expectedCode, clientErr.ErrorCode())
		}

	})

	// Test case 2: HTTP non-200 status code
	t.Run("HTTP non-200 status code", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized) // 401
			w.Write([]byte(`{"error": "Invalid API Key"}`))
		}))
		defer server.Close()

		client := NewAlphaVantageClient(server.URL, "invalid-key")
		stocks, err := client.FetchDaily("IBM")

		if stocks != nil {
			t.Error("Expected nil stocks, but got some data")
		}

		expectedError := &client_errors.HTTPStatusError{}
		expectedCode := expectedError.ErrorCode()

		// Check specific error type
		if !errors.As(err, &expectedError) {
			t.Errorf("Expected error of type *client_errors.HTTPStatusError, but got %T", err)
		}

		// Check interface implementation and error code
		clientErr, ok := err.(client_errors.ClientError)
		if !ok {
			t.Errorf("Error does not implement ClientError interface")
		} else if clientErr.ErrorCode() != expectedCode {
			t.Errorf("Expected error code %s, but got %s", expectedCode, clientErr.ErrorCode())
		}

	})
}
