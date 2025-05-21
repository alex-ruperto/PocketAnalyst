package controllers

import (
	"PocketAnalyst/services"
	"encoding/json"
	"log"
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

// HandleStockFetchRequest handles requests to fetch and store new stock data
func (sc *StockController) HandleStockFetchRequest(w http.ResponseWriter, r *http.Request) {
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
	count, err := sc.stockService.SynchronizeStockData(r.Context(), symbol)
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

func (sc *StockController) HandleStockHistoryRequest(w http.ResponseWriter, r *http.Request) {
	// Parse request parameters
	symbol := r.URL.Query().Get("symbol")
	if symbol == "" {
		http.Error(w, "Symbol parameter is required", http.StatusBadRequest)
		return
	}

	// Default to last 30 days if not provided
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -30)

	// Parse date ranges if provided in query parameters
	if startDateStr := r.URL.Query().Get("start_date"); startDateStr != "" {
		if parsedDate, err := time.Parse("2006-01-02", startDateStr); err == nil {
			startDate = parsedDate
		} else {
			log.Printf("Invalid start date format: %s", startDateStr)
		}
	}

	if endDateStr := r.URL.Query().Get("end_date"); endDateStr != "" {
		if parsedDate, err := time.Parse("2006-01-02", endDateStr); err == nil {
			endDate = parsedDate
		} else {
			log.Printf("Invalid end date format: %s", endDateStr)
		}
	}

	// Get stock history from service layer
	stocks, err := sc.stockService.GetStockHistory(r.Context(), symbol, startDate, endDate)
	if err != nil {
		http.Error(w, "Error retrieving stock data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return data as JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(stocks); err != nil {
		http.Error(w, "Error encoding response: "+err.Error(), http.StatusInternalServerError)
	}
}
