package errors

import (
	"fmt"
)

// ModelValidationError represents model validation errors
type ModelValidationError struct {
	Model   string
	Field   string
	Message string
}

// Error returns a formatted error message for ModelValidationError.
func (e *ModelValidationError) Error() string {
	return fmt.Sprintf("%s validation failed for %s: %s",
		e.Model,
		e.Field,
		e.Message)
}

// NewModeValidationError creates a new model validation error
// given the model, field, and the message.
func NewModelValidationError(model, field, message string) error {
	return &ModelValidationError{
		Model:   model,
		Field:   field,
		Message: message,
	}
}
