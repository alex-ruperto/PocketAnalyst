package models

import (
	"fmt"
	"pocketanalyst/pkg/errors"
	"time"
)

// Stock represents historical stock price data in the database.
type Stock struct {
	PriceID          int       `json:"price_id"` // SERIAL, auto-incrementing PK.
	CompanyID        int       `json:"company_id"`
	Symbol           string    `json:"symbol"`
	Date             time.Time `json:"date"`
	OpenPrice        float64   `json:"open_price"`
	HighPrice        float64   `json:"high_price"`
	LowPrice         float64   `json:"low_price"`
	ClosePrice       float64   `json:"close_price"`
	AdjustedClose    float64   `json:"adjusted_close"`
	Volume           float64   `json:"volume"`
	DividendAmount   float64   `json:"dividend_amount"`
	SplitCoefficient float64   `json:"split_coefficient"`
	DataSource       string    `json:"data_source"`
	LastUpdated      time.Time `json:"last_updated"`
}

// String implements the Stringer interface for Stock
func (s *Stock) String() string {
	return fmt.Sprintf(
		"Stock{Symbol: %s, Date: %s, Open: %.2f, High: %.2f, Low: %.2f, Close: %.2f, Volume: %.0f}",
		s.Symbol,
		s.Date.Format("2006-01-02"),
		s.OpenPrice,
		s.HighPrice,
		s.LowPrice,
		s.ClosePrice,
		s.Volume,
	)
}

// Validate checks if the stock data meets all logical rules
func (s *Stock) Validate() error {
	switch {
	case s.Symbol == "":
		return errors.NewModelValidationError("Stock", "symbol", "symbol is required")
	case s.CompanyID < 0:
		return errors.NewModelValidationError("Stock", "company_id", "company_id must be a positive number")
	case s.Date.IsZero():
		return errors.NewModelValidationError("Stock", "date", "date cannot be empty")
	case s.OpenPrice < 0:
		return errors.NewModelValidationError("Stock", "open_price", "open_price cannot be negative")
	case s.HighPrice < 0:
		return errors.NewModelValidationError("Stock", "high_price", "high_price cannot be negative")
	case s.LowPrice < 0:
		return errors.NewModelValidationError("Stock", "low_price", "low_price cannot be negative")
	case s.ClosePrice < 0:
		return errors.NewModelValidationError("Stock", "close_price", "close_price cannot be negative")
	case s.AdjustedClose < 0:
		return errors.NewModelValidationError("Stock", "adjusted_close_price", "adjusted_close_price cannot be negative")
	case s.HighPrice < s.LowPrice && s.HighPrice > 0 && s.LowPrice > 0:
		return errors.NewModelValidationError("Stock", "price_range", "high price cannot be less than low price")
	case s.Volume < 0:
		return errors.NewModelValidationError("Stock", "volume", "volume cannot be negative")
	}
	return nil
}
