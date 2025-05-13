package clients

import (
	"http"
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
	tests := []struct {
		name           string
		serverResponse func(w http.Response, r *http.Request)
		expectedError  string
	}{
		{
			name:           "HTTP Request Failure",
			serverResponse: nil,
			expectedError:  "failed to make a request to Alpha Vantage",
		},
		{
			name:		"Non-200 Status Code",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusUnauthorized) // 401
				w.Write([]byte(`{"error: "Invalid API Key}`))
			},
			expectedError:	"Alpha Vantage API returned status code 401",
		}
	}

}
