package models

import (
	"errors"
	"time"
)

// Stock represent a stock data entry in the Postgres DB.
type Stock struct {
	StockDataID   int       `json:"stock_data_id"`
	CompanyID     int       `json:"company_id"`
	Date          time.Time `json:"date"`
	OpenPrice     float64   `json:"open_price"`
	HighPrice     float64   `json:"high_price"`
	LowPrice      float64   `json:"low_price"`
	ClosePrice    float64   `json:"close_price"`
	AdjClosePrice float64   `json:"adj_close_price"`
	Volume        float64   `json:"volume"`
}

// Validate ensures the stock data meets logical rules
func (s *Stock) Validate() error {
	// This tagless switch will evaluate all the following conditions
	switch {
	case s.CompanyID <= 0:
		return errors.New("company_id must be a positive integer")

	case s.OpenPrice < 0 || s.HighPrice < 0 || s.LowPrice < 0 || s.ClosePrice < 0 || s.AdjClosePrice < 0:
		return errors.New("price values cannot be negative")

	case s.HighPrice < s.LowPrice:
		return errors.New("high price cannot be less than low price")

	case s.Volume < 0:
		return errors.New("volume cannot be negative")

	case s.Date.IsZero():
		return errors.New("date cannot be empty")
	}

	// No errors found
	return nil
}
