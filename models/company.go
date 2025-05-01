package models

import (
	"errors"
	"time"
)

// Company represents a company in the database
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

func (c *Company) Validate() error {
	switch {
	case c.Symbol == "":
		return errors.New("symbol is required")

	case c.Name == "":
		return errors.New("name is required")
	}
	return nil
}
