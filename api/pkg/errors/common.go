package errors

import (
	"errors"
	"fmt"
)

// Standard errors have utility functions for error wrapping
// Allows consumers to use this package's errors exclusively
// without also importing the standard errors package.
var (
	// Reports whether any error in err's tree matches target.
	Is = errors.Is

	// Finds the first error in err's tree that matches target,
	// and if one is found, target to that error value and
	// return true.
	As = errors.As

	// Returns the result of calling the Unwrap method on err,
	// if err's type contains an Unwrap method returning error.
	// Otherwise, Unwrap returns nil.
	Unwrap = errors.Unwrap
)

// NotFoundError occurs when a requested resource is not found.
type NotFoundError struct {
	EntityType string
	ID         any
}

// Implements the error interface, returning a formatted NotFoundError.
func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s with ID %s was not found",
		e.EntityType,
		e.ID)
}

// NewNotFoundError creates a new not found error for the given entity
// type and ID.
func NewNotFoundError(entityType string, id any) error {
	return &NotFoundError{
		EntityType: entityType,
		ID:         id,
	}
}

// ServiceError wraps an error that occurred during a service operation
type ServiceError struct {
	Operation string
	Err       error
}

// Error returns a formatted error message for ServiceError.
func (e *ServiceError) Error() string {
	return fmt.Sprintf("service error during %s: %v",
		e.Operation,
		e.Err)
}

// Unwrap returns the wrapped error.
func (e *ServiceError) Unwrap() error {
	return e.Err
}

// NewServiceError creates a new ServiceError that wraps the given error.
func NewServiceError(operation string, err error) error {
	return &ServiceError{
		Operation: operation,
		Err:       err,
	}
}
