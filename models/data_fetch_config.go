package models

import (
	"errors"
	"time"
)

// DataFetchConfig represents the configuration for external API data fetching
type DataFetchConfig struct {
	ConfigID    int       `json:"config_id"`
	Symbol      string    `json:"symbol"`
	DataSource  string    `json:"data_source"`
	APIKey      string    `json:"api_key"`
	StartDate   time.Time `json:"start_date,omitempty"`
	EndDate     time.Time `json:"end_date,omitempty"`
	Frequency   string    `json:"frequency"` // "daily", "weekly", "monthly"
	LastFetched time.Time `json:"last_fetched"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
}

// DataSource constants
const (
	SourceAlphaVantage = "AlphaVantage"
	SourceYahooFinance = "YahooFinance"
	SourceIEXCloud     = "IEXCloud"
	SourcePolygon      = "Polygon"
)

// Frequencies constants
const (
	FrequencyIntraday = "intraday"
	FrequencyDaily    = "daily"
	FrequencyWeekly   = "weekly"
	FrequencyMonthly  = "monthly"
)

// Validate ensures the data fetch config meets all logical rules
func (c *DataFetchConfig) Validate() error {
	switch {
	case c.Symbol == "":
		return errors.New("symbol is required")

	case c.DataSource == "":
		return errors.New("data source is required")

	case c.Frequency == "":
		return errors.New("frequency is required")
	}

	// Validate specific frequencies
	validFrequencies := map[string]bool{
		"daily":    true,
		"weekly":   true,
		"monthly":  true,
		"intraday": true,
	}

	if !validFrequencies[c.Frequency] {
		return errors.New("invalid frequency: must be daily, weekly, monthly, or intraday")
	}

	// Validate date range if both dates are provided
	if !c.StartDate.IsZero() && !c.EndDate.IsZero() {
		if c.EndDate.Before(c.StartDate) {
			return errors.New("end date cannot be before start date")
		}
	}

	// Validate data sources
	validSources := map[string]bool{
		"AlphaVantage": true,
		"YahooFinance": true,
		"IEXCloud":     true,
		"Polygon":      true,
	}

	if !validSources[c.DataSource] {
		return errors.New("invalid data source")
	}

	return nil
}

// ShouldFetchToday checks if this config should be fetched today
// based on last fetch date and frequency
func (c *DataFetchConfig) ShouldFetchToday() bool {
	if !c.IsActive {
		return false
	}

	// If never fetched before, fetch now
	if c.LastFetched.IsZero() {
		return true
	}

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	lastFetchDate := time.Date(
		c.LastFetched.Year(),
		c.LastFetched.Month(),
		c.LastFetched.Day(),
		0, 0, 0, 0,
		c.LastFetched.Location(),
	)

	// Don't fetch more than once per day
	if today.Equal(lastFetchDate) {
		return false
	}

	// If it's been a day or more since last fetch, fetch again
	return today.After(lastFetchDate)
}

// NewDataFetchConfig creates a new DataFetchConfig with default values
func NewDataFetchConfig(symbol string, dataSource string) *DataFetchConfig {
	return &DataFetchConfig{
		Symbol:     symbol,
		DataSource: dataSource,
		Frequency:  FrequencyDaily,
		IsActive:   true,
		CreatedAt:  time.Now(),
	}
}
