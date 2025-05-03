package repositories

import (
	"GoSimpleRest/models"
	"context"
	"database/sql"
	"fmt"
	"time"
)

// StockRepository handles database operations for stocks
type StockRepository struct {
	db *sql.DB
}

// NewStockRepository creates a new stock repository
func NewStockRepository(db *sql.DB) *StockRepository {
	return &StockRepository{db: db}
}
