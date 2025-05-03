package models

import (
	"GoSimpleREST/errors"
	"time"
)

// Company represents a company entity in the database
type Company struct {
	CompanyID int       `json:"company_id"`
	Symbol    string    `json:"symbol"`
	Name      string    `json:"name"`
	Sector    string    `json:"sector"`
	Industry  string    `json:"industry"`
	Exchange  string    `json:"exchange"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

// Validate ensures the stock data meets all logical rules
func (c *Company) Validate() error {
	switch {
	case c.Symbol == "":
		return errors.NewModelValidationError("Company", "symbol", "symbol is required")
	case c.Name == "":
		return errors.NewModelValidationError("Company", "name", "name is required.")
	}
	return nil
}
