package repositories

import (
	"PocketAnalyst/models"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
)

// DataSourceRepository handles DB operations for data sources
type DataSourceRepository struct {
	db *sql.DB
}

// NewDataSourceRepository creates a new data source repository
func NewDataSourceRepository(db *sql.DB) *DataSourceRepository {
	return &DataSourceRepository{db: db}
}

// GetByName retrieves a data source by it name
func (dsr *DataSourceRepository) GetByName(ctx context.Context, name string) (*models.DataSource, error) {
	query := `
		SELECT *
		FROM data_sources
		WHERE source_name = $1
	`
	var ds models.DataSource
	var configJSON []byte

	err := dsr.db.QueryRowContext(ctx, query, name).Scan(
		&ds.SourceID,
		&ds.SourceName,
		&ds.SourceType,
		&ds.BaseURL,
		&ds.RateLimitPerMinute,
		&ds.RateLimitPerDay,
		&configJSON,
		&ds.IsActive,
		&ds.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("data source with name '%s' not found", name)
		}
		return nil, fmt.Errorf("Error retrieving data source: %w", err)
	}

	// Parse config parameters
	if len(configJSON) > 0 {
		if err := json.Unmarshal(configJSON, &ds.ConfigParameters); err != nil {
			return nil, fmt.Errorf("error parsing config parameters: %w", err)
		}
	} else {
		ds.ConfigParameters = make(map[string]any)
	}

	return &ds, nil
}
