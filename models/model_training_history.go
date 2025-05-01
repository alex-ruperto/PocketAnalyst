package models

import (
	"errors"
	"time"
)

// ModelTrainingHistory represents a record of model training
type ModelTrainingHistory struct {
	TrainingID           int       `json:"training_id"`
	ModelVersion         string    `json:"model_version"`
	TrainingDate         time.Time `json:"training_date"`
	ModelType            string    `json:"model_type"`
	Parameters           string    `json:"parameters"` // JSON string of parameters
	TrainingAccuracy     float64   `json:"training_accuracy"`
	ValidationAccuracy   float64   `json:"validation_accuracy"`
	FeaturesUsed         string    `json:"features_used"`
	TrainingDurationSecs int       `json:"training_duration_seconds"`
	Notes                string    `json:"notes"`
}

// Validate ensures the training history data meets all business rules
func (m *ModelTrainingHistory) Validate() error {
	switch {
	case m.ModelVersion == "":
		return errors.New("model version is required")

	case m.ModelType == "":
		return errors.New("model type is required")

	case m.TrainingAccuracy < 0 || m.TrainingAccuracy > 1:
		return errors.New("training accuracy must be between 0 and 1")

	case m.ValidationAccuracy < 0 || m.ValidationAccuracy > 1:
		return errors.New("validation accuracy must be between 0 and 1")
	}

	return nil
}

const (
	ModelLSTM         = "LSTM"          // Long Short-Term Memory
	ModelGRU          = "GRU"           // Gated Recurrent Unit
	ModelRandomForest = "RANDOM_FOREST" // Random Forest
	ModelXGBoost      = "XGBOOST"       // XGBoost
	ModelSVM          = "SVM"           // Support Vector Machine
	ModelLinearReg    = "LINEAR_REG"    // Linear Regression
	ModelEnsemble     = "ENSEMBLE"      // Ensemble Model
)
