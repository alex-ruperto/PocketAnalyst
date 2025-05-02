package models

import (
	"GoSimpleREST/errors"
	"time"
)

// Stock represents historical stock price data in the database.
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
