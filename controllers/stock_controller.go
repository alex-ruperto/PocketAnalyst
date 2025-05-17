package controllers

import (
	"PocketAnalyst/services"
	"encoding/json"
	"net/http"
	"time"
)

// StockController handles HTTP Requests related to stocks
type StockController struct {
	stockService *services.StockService
}

// NewStockController creates a new instance of StockController
func NewStockController(stockService *services.StockService) *StockController {
	return &StockController{
		stockService: stockService,
	}
}

// FetchStockData handles requests to fetch and store new stock data
func (sc *StockController) FetchStockData(w http.ResponseWriter, r *http.Request)
