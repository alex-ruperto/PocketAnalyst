package clients

import (
	"net/http"
	"net/http/httptest"
	"os"
	"pocketanalyst/pkg/errors/client_errors"
	// "reflect"
	"testing"
)

func TestAlphaVantageClient_FetchDailyPricesFromAPI(t *testing.T) {
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
	stocks, err := client.FetchDailyPricesFromAPI("AAPL")
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
		stocks, err := client.FetchDailyPricesFromAPI("IBM")

		if stocks != nil {
			t.Error("Expected nil stocks, but got some data.")
		}

		expectedErr := &client_errors.HTTPRequestError{}
		errorValidate(t, err, expectedErr)

	})

	// Test case 2: HTTP non-200 status code
	t.Run("HTTP non-200 status code", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized) // 401
			w.Write([]byte(`{"error": "Invalid API Key"}`))
		}))
		defer server.Close()

		client := NewAlphaVantageClient(server.URL, "invalid-key")
		stocks, err := client.FetchDailyPricesFromAPI("IBM")

		if stocks != nil {
			t.Error("Expected nil stocks, but got some data")
		}

		expectedErr := &client_errors.HTTPStatusError{}
		errorValidate(t, err, expectedErr)
	})

	// Test case 3: API Error
	t.Run("API Error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"Error Message": "Invalid API call"}`))
		}))
		defer server.Close()

		client := NewAlphaVantageClient(server.URL, "none")
		stocks, err := client.FetchDailyPricesFromAPI("IBM")
		if stocks != nil {
			t.Error("Expected nil, but got some data")
		}

		expectedErr := &client_errors.APIError{}
		errorValidate(t, err, expectedErr)
	})
}

// Validates any error type that implements the ClientError interface.
func errorValidate(t *testing.T, err error, expectedErr client_errors.ClientError) {
	if err == nil {
		t.Error("Expected error, but got nil")
		return
	}

	// Check if the error is a client error
	clientErr, ok := err.(client_errors.ClientError)
	if !ok {
		t.Errorf("Error does not implement the ClientError interface")
		return
	}

	// Not strictly necessary because we are already comparing error codes, eliminating the need for reflection
	// Check type equality
	// expectedType := reflect.TypeOf(expectedErr)
	// actualType := reflect.TypeOf(err)
	// if expectedType != actualType {
	// 	t.Errorf("Expected error of type %v, but got %v", expectedType, actualType)
	//	return
	// }

	// Check error code equality
	expectedCode := expectedErr.ErrorCode()
	if clientErr.ErrorCode() != expectedCode {
		t.Errorf("Expected error code %s, but got %s", expectedCode, clientErr.ErrorCode())
	}
}
