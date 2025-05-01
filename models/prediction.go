package models

import (
	"errors"
	"time"
)

// StockPrediction represents ML model predictions for a stock

type StockPrediction struct {
	PredictionID   int       `json:"prediction_id"`
	Symbol         string    `json:"symbol"`
	TargetDate     time.Time `json:"target_date"`     // Date for which prediction was made
	PredictionDate time.Time `json:"prediction_date"` // When the prediction was made
	PredictionType string    `json:"prediction_type"` // E.g., "DIRECTION", "PRICE", "VOLATILITY"
	PredictedValue float64   `json:"predicted_value"` // Predicted value (direction: -1 to 1, price: actual value)
	Confidence     float64   `json:"confidence"`      // Model confidence 0-1
	ModelVersion   string    `json:"model_version"`   // Version of ML model used
	FeaturesUsed   string    `json:"features_used"`   // Description of features used
	ActualValue    float64   `json:"actual_value"`    // Actual value once known
	AccuracyMetric float64   `json:"accuracy_metric"` // How accurate the prediction was
}

// PredictionType constraints
const (
	PredictionDirection  = "DIRECTION"
	PredictionPrice      = "PRICE"
	PredictionVolatility = "VOLATILITY"
	PredictionReturn     = "RETURN"
)

// Validate ensures the prediction data meets all logical rules
func (p *StockPrediction) Validate() error {
	switch {
	case p.Symbol == "":
		return errors.New("symbol is required")

	case p.TargetDate.IsZero():
		return errors.New("target date cannot be empty")

	case p.PredictionType == "":
		return errors.New("prediction type is required")

	case p.ModelVersion == "":
		return errors.New("model version is required")

	case p.Confidence < 0 || p.Confidence > 1:
		return errors.New("confidence must be between 0 and 1")
	}

	// Validate specific prediction types
	switch p.PredictionType {
	case "DIRECTION":
		if p.PredictedValue < -1 || p.PredictedValue > 1 {
			return errors.New("direction prediction must be between -1 and 1")
		}
	case "PRICE":
		if p.PredictedValue < 0 {
			return errors.New("price prediction cannot be negative")
		}
	case "VOLATILITY":
		if p.PredictedValue < 0 {
			return errors.New("volatility prediction cannot be negative")
		}
	}

	return nil
}
