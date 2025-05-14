package models

import (
	"PocketAnalyst/errors"
	"time"
)

// DataSource represents a data source configuration
type DataSource struct {
	SourceID           int            `json:"source_id"`
	SourceName         string         `json:"source_name"`
	SourceType         string         `json:"source_type"`
	BaseURL            string         `json:"base_url"`
	RateLimitPerMinute string         `json:"rate_limit_per_minute"`
	RateLimitPerDay    string         `json:"rate_limit_per_day"`
	ConfigParameters   map[string]any `json:"config_parameters"`
	IsActive           bool           `json:"is_active"`
	CreatedAt          time.Time      `json:"created_at"`
}

// Validate ensures the data source meets all logical requirements
func (ds *DataSource) Validate() error {
	switch {
	case ds.SourceName == "":
		return errors.NewModelValidationError("DataSource", "source_name", "source_name cannot be empty")
	case ds.SourceType == "":
		return errors.NewModelValidationError("DataSource", "source_type", "source_type cannot be empty")
	}
	return nil
}
