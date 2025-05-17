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
func (sc *StockController) FetchStockData(w http.ResponseWriter, r *http.Request) {
	// Only allow POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse requset parameters
	symbol := r.URL.Query().Get("symbol")
	if symbol == "" {
		http.Error(w, "symbol parameter is required", http.StatusBadRequest)
		return
	}

	// Fetch and store stock data in DB
	count, err := sc.stockService.FetchAndStoreStockData(r.Context(), symbol)
	if err != nil {
		http.Error(w, "Error fetching stock data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return success response
	response := map[string]any{
		"success":           true,
		"records_processed": count,
		"message":           "Successfully fetched and stored stock data",
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding response: "+err.Error(), http.StatusInternalServerError)
	}
}
