package clients

import (
	"os"
	"testing"
	"time"
)

// TestFMPClient_FetchDaily_Integration tests the FMP client with real API calls
// This test makes actual HTTP requests to the FMP API
func TestFMPClient_FetchDaily_Integration(t *testing.T) {
	// Get the API Key from environment variable
	apiKey := os.Getenv("FMP_API_KEY")
	if apiKey == "" {
		// Fall back to the key you provided for testing
		apiKey = "KLf66Ai0X9ko34vNwZCgwc1qitU4Wowk"
		t.Logf("Using hardcoded API key for testing")
	}

	// Create a real FMP client with the actual base URL
	client := NewFMPClient("https://financialmodelingprep.com", apiKey)

	// Test with NVDA symbol (matches your example)
	symbol := "NVDA"
	t.Logf("Testing FMP client with symbol: %s", symbol)

	// Make the actual API call
	stocks, err := client.FetchDaily(symbol)
	if err != nil {
		t.Fatalf("Failed to fetch data from FMP: %v", err)
	}

	// Verify we got some data back
	if len(stocks) == 0 {
		t.Error("No stock data was returned from FMP")
		return
	}

	t.Logf("Successfully fetched %d stock records from FMP", len(stocks))

	// Verify the first stock record has expected data
	firstStock := stocks[0]

	// Basic validation checks
	if firstStock.Symbol != symbol {
		t.Errorf("Expected symbol %s, got %s", symbol, firstStock.Symbol)
	}

	if firstStock.Date.IsZero() {
		t.Error("Stock date should not be zero")
	}

	if firstStock.OpenPrice <= 0 {
		t.Errorf("Open price should be positive, got %f", firstStock.OpenPrice)
	}

	if firstStock.HighPrice <= 0 {
		t.Errorf("High price should be positive, got %f", firstStock.HighPrice)
	}

	if firstStock.LowPrice <= 0 {
		t.Errorf("Low price should be positive, got %f", firstStock.LowPrice)
	}

	if firstStock.ClosePrice <= 0 {
		t.Errorf("Close price should be positive, got %f", firstStock.ClosePrice)
	}

	if firstStock.Volume <= 0 {
		t.Errorf("Volume should be positive, got %f", firstStock.Volume)
	}

	// Verify price relationships make sense
	if firstStock.HighPrice < firstStock.LowPrice {
		t.Errorf("High price (%f) should be >= low price (%f)",
			firstStock.HighPrice, firstStock.LowPrice)
	}

	// Verify data source is set correctly
	if firstStock.DataSource != "FMP" {
		t.Errorf("Expected data source 'FMP', got '%s'", firstStock.DataSource)
	}

	// Verify the date is recent (within last 2 years for active trading)
	now := time.Now()
	twoYearsAgo := now.AddDate(-2, 0, 0)
	if firstStock.Date.Before(twoYearsAgo) {
		t.Errorf("Stock date seems too old: %s", firstStock.Date.Format("2006-01-02"))
	}

	// Log detailed information about the first few stocks
	t.Log("First few stock data points:")
	for i, stock := range stocks {
		if i >= 3 { // Only show first 3 for brevity
			break
		}
		t.Logf("Stock #%d: %s", i+1, stock.String())

		// Additional detailed logging
		t.Logf("  Date: %s", stock.Date.Format("2006-01-02"))
		t.Logf("  OHLC: %.2f / %.2f / %.2f / %.2f",
			stock.OpenPrice, stock.HighPrice, stock.LowPrice, stock.ClosePrice)
		t.Logf("  Volume: %.0f", stock.Volume)
		t.Logf("  Data Source: %s", stock.DataSource)
		t.Logf("  Last Updated: %s", stock.LastUpdated.Format("2006-01-02 15:04:05"))
	}
}
