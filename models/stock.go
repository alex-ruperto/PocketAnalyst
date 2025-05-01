package models

import (
	"errors"
	"time"
)

// Stock represent a stock data entry in the Postgres DB.
type Stock struct {
	StockDataID        int       `json:"stock_data_id"`
	CompanyID          int       `json:"company_id"`
	Symbol             string    `json:"symbol"`
	Date               time.Time `json:"date"`
	OpenPrice          float64   `json:"open_price"`
	HighPrice          float64   `json:"high_price"`
	LowPrice           float64   `json:"low_price"`
	ClosePrice         float64   `json:"close_price"`
	AdjustedClosePrice float64   `json:"adjusted_close_price"`
	Volume             float64   `json:"volume"`
	DividendAmount     float64   `json:"dividend_amount"`
	SplitCoefficient   float64   `json:"split_coefficient"`
	DataSource         string    `json:"data_source"`
	LastUpdated        time.Time `json:"last_updated"`
}

// Validate ensures the stock data meets logical rules
func (s *Stock) Validate() error {
	// This tagless switch will evaluate all the following conditions
	switch {
	case s.Symbol == "":
		return errors.New("symbol is required")

	case s.CompanyID <= 0:
		return errors.New("company_id must be a positive integer")

	case s.Date.IsZero():
		return errors.New("date cannot be empty")

	case s.OpenPrice < 0 || s.HighPrice < 0 || s.LowPrice < 0 || s.ClosePrice < 0 || s.AdjustedClosePrice < 0:
		return errors.New("price values cannot be negative")

	case s.HighPrice < s.LowPrice && s.HighPrice > 0 && s.LowPrice > 0:
		return errors.New("high price cannot be less than low price")

	case s.Volume < 0:
		return errors.New("volume cannot be negative")

	case s.SplitCoefficient < 0:
		return errors.New("split coefficient cannot be negative")

	case s.DividendAmount < 0:
		return errors.New("dividend amount cannot be negative")
	}

	// No errors found
	return nil
}
