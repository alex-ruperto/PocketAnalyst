package clients

import (
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
