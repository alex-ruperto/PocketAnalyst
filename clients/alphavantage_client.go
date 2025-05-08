package clients

import (
	"GoSimpleREST/models"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"
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

// Fetch daily adjusted stock prices from Alpha Vantage
// Alpha Vantage API documentation https://www.alphavantage.co/documentation/
func (avc *AlphaVantageClient) FetchDailyAdjusted(symbol string) ([]*models.Stock, error) {
	// Construct URL with required parameters
	// TIME_SERIES_DAILY_ADJUSTED returns daily adjusted time series
	// outputsize=compact returns the latest 100 data points
	url := fmt.Sprintf("%s?function=TIME_SERIES_DAILY_ADJUSTED&symbol=%s&outputsize=compact&apikey=%s",
		avc.BaseURL, symbol, avc.APIKey)

	// Make HTTP request.
	resp, err := avc.Client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to make a request to Alpha Vantage %w", err)
	}
	defer resp.Body.Close() // Ensure the response body is closed to prevent resource leaks.

	// Check for successful HTTP response
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Alpha Vantage API returned status code %d", resp.StatusCode)
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to read Alpha Vantage response: %w", err)
	}

	// Parse the JSON response
	var response map[string]any
	// Unmarshal only needs to read body. We need to modify the original response so we pass response by reference here.
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse Alpha Vantage response: %w", err)
	}

	// Check for API error messages
	if errorMsg, exists := response["Error Message"]; exists {
		return nil, fmt.Errorf("Alpha Vantage API error: %v", errorMsg)
	}

	// Extract the time series data
	timeSeriesKey := "Time Series (Daily)"
	timeSeries, ok := response[timeSeriesKey].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("could not find time series data in response.")
	}

	// Converts the data into Stock models
	stocks := make([]*models.Stock, 0, len(timeSeries))

	// Sort dates to ensure we process them in chronological order
	dates := make([]string, 0, len(timeSeries))
	for date := range timeSeries {
		dates = append(dates, date)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(dates))) // Sort in descending order (newest first)

	// Process the first 5 entries for display
	count := 0
	for _, dateStr := range dates {
		if count >= 5 {
			break // Only process the first five entries
		}

		dailyData, ok := timeSeries[dateStr].(map[string]any)
		if !ok {
			continue
		}

		// Parse the date
		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			continue // Skip dates we can't parse
		}

		stock := &models.Stock{
			Symbol:           symbol,
			Date:             date,
			OpenPrice:        parseFloat(dailyData, "1. open"),
			HighPrice:        parseFloat(dailyData, "2. high"),
			LowPrice:         parseFloat(dailyData, "3. low"),
			ClosePrice:       parseFloat(dailyData, "4. close"),
			AdjustedClose:    parseFloat(dailyData, "5. adjusted close"),
			Volume:           parseFloat(dailyData, "6. volume"),
			DividendAmount:   parseFloat(dailyData, "7. dividend amount"),
			SplitCoefficient: parseFloat(dailyData, "8. split coefficient"),
			DataSource:       "AlphaVantage",
			LastUpdated:      time.Now(),
		}
	}

}

func parseFloat(data map[string]any, key string) float64 {
	if val, ok := data[key].(string); ok {
		if f, err := strconv.ParseFloat(val, 64); err == nil {
			return f
		}
	}
	return 0.0
}
