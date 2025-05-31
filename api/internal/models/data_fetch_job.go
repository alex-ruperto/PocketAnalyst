package models

import (
	"pocketanalyst/pkg/errors"
	"time"
)

// DataFetchJob represents a scheduled data fetching job
type DataFetchJob struct {
	JobID         int            `json:"job_id"`
	SourceID      int            `json:"source_id"`
	EntityType    string         `json:"entity_type"`
	EntityValue   string         `json:"entity_value"`
	DataType      string         `json:"data_type"`
	Frequency     string         `json:"frequency"`
	Parameters    map[string]any `json:"parameters"`
	LastExecution time.Time      `json:"last_execution"`
	LastSuccess   time.Time      `json:"last_success"`
	NextScheduled time.Time      `json:"next_scheduled"`
	Status        string         `json:"status"`
	IsActive      string         `json:"is_active"`
	CreatedAt     string         `json:"created_at"`
}

// Validate ensures the data fetch job meets all logical requirements
func (j *DataFetchJob) Validate() error {
	switch {
	case j.SourceID <= 0:
		return errors.NewModelValidationError("DataFetchJob", "source_id", "source_id must be positive")
	case j.EntityType == "":
		return errors.NewModelValidationError("DataFetchJob", "entity_type", "entity_type is required")
	case j.EntityValue == "":
		return errors.NewModelValidationError("DataFetchJob", "entity_value", "entity_value is required")
	case j.DataType == "":
		return errors.NewModelValidationError("DataFetchJob", "data_type", "data_type is required")
	case j.Frequency == "":
		return errors.NewModelValidationError("DataFetchJob", "frequency", "frequency is required")
	case j.NextScheduled.IsZero():
		return errors.NewModelValidationError("DataFetchJob", "next_scheduled", "next scheduled time is required")
	}
	return nil
}
