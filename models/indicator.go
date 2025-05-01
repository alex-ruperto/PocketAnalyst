package models

import (
	"errors"
	"time"
)

// StockIndicator represents calculated technical indicators for a stock
type StockIndicator struct {
	IndicatorID   int       `json:"indicator_id"`
	StockDataID   int       `json:"stock_data_id"`
	Symbol        string    `json:"symbol"`
	Date          time.Time `json:"date"`
	IndicatorType string    `json:"indicator_type"`       // E.g., SMA, BOLLINGER, RSI, EMA
	Period        int       `json:"period"`               // Time period for the indicator (e.g., 14-day RSI)
	Value         float64   `json:"value"`                // Primary indicator value
	UpperBand     float64   `json:"upper_band,omitempty"` // For indicators with bands like Bollinger
	LowerBand     float64   `json:"lower_band,omitempty"` // For indicators with bands like Bollinger
	CreatedAt     time.Time `json:"created_at"`
}

// Indicator type constraints
const (
	IndicatorSMA       = "SMA"       // Simple Moving Average
	IndicatorEMA       = "EMA"       // Exponential Moving Average
	IndicatorRSI       = "RSI"       // Relative Strength Index
	IndicatorMACD      = "MACD"      // Moving Average Convergence Divergence
	IndicatorBollinger = "BOLLINGER" // Bollinger Bands
	IndicatorATR       = "ATR"       // Average True Range
	IndicatorStochRSI  = "STOCHRSI"  // Stochastic RSI
	IndicatorADX       = "ADX"       // Average Directional Index
	IndicatorOBV       = "OBV"       // On-Balance Volume
)

// Ensures the indicator data meets all logical rules
func (i *StockIndicator) Validate() error {
	switch {
	case i.Symbol == "":
		return errors.New("symbol is required")

	case i.Date.IsZero():
		return errors.New("date cannot be empty")

	case i.IndicatorType == "":
		return errors.New("indicator type is required")

	case i.Period <= 0:
		return errors.New("period must be a positive integer")
	}

	// Validate specific indicator types
	switch i.IndicatorType {
	case "RSI":
		if i.Value < 0 || i.Value > 100 {
			return errors.New("RSI value must be between 0 and 100")
		}

	case "BOLLINGER", "KELTNER", "DONCHIAN":
		if i.UpperBand == 0 || i.LowerBand == 0 {
			return errors.New("band indicators must have upper and lower band values")
		}
		if i.UpperBand < i.LowerBand {
			return errors.New("upper band value cannot be less than lower band value")
		}
	}
	return nil
}

// NewStockIndicator creates a new stock indicator with default values
func NewStockIndicator(symbol string, date time.Time, indicatorType string, period int) *StockIndicator {
	return &StockIndicator{
		Symbol:        symbol,
		Date:          date,
		IndicatorType: indicatorType,
		Period:        period,
		CreatedAt:     time.Now(),
	}
}
