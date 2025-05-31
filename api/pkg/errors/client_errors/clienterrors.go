// Package client_erorrs handles client-related errors
//
//	ClientError: Common interface for all API client errors.
package client_errors

import (
	"fmt"
)

// Base interface for all API client errors
type ClientError interface {
	error
	ErrorCode() string
}

// HTTPRequestError occurs when a request to an external API fails.
// E.g., network issues, DNS problems
type HTTPRequestError struct {
	URL        string
	InnerError error
}

func (e *HTTPRequestError) Error() string {
	return fmt.Sprintf("failed to make a request to %s: %v", e.URL, e.InnerError)
}

func (e *HTTPRequestError) ErrorCode() string { // Implicit implementation of Error interface
	return "HTTP_REQUEST_ERROR"
}

func (e *HTTPRequestError) Unwrap() error {
	return e.InnerError
}

// Returns a new HTTPRequestError.
//
//	HTTPRequestError: Occurs when a request to an external API Fails.
//	url: The request URL
//	err: The inner error
func NewHTTPRequestError(url string, err error) *HTTPRequestError {
	return &HTTPRequestError{
		URL:        url,
		InnerError: err,
	}
}

// HTTPStatusError occurs when a request does NOT return a 200 status code.
type HTTPStatusError struct {
	URL          string
	StatusCode   int
	ResponseBody string
}

func (e *HTTPStatusError) Error() string {
	return fmt.Sprintf("API returned status code %d for URL %s. Response Body: %s", e.StatusCode, e.URL, e.ResponseBody)
}

func (e *HTTPStatusError) ErrorCode() string {
	return "HTTP_STATUS_ERROR"
}

// Returns a new HTTPStatusError
//
//	HTTPStatusError: Occurs when the returned status code is not 200.
//	url: The request URL
//	statusCode: The status code returned from the HTTP request
//	responseBody: The response body from the request
func NewHTTPStatusError(url string, statusCode int, responseBody string) *HTTPStatusError {
	return &HTTPStatusError{
		URL:          url,
		StatusCode:   statusCode,
		ResponseBody: responseBody,
	}
}

// ResponseReadError occurs when reading the API response body fails.
type ResponseReadError struct {
	InnerError error
}

func (e *ResponseReadError) Error() string {
	return fmt.Sprintf("failed to read response body: %v", e.InnerError)
}

func (e *ResponseReadError) ErrorCode() string {
	return "RESPONSE_READ_ERROR"
}

func (e *ResponseReadError) Unwrap() error {
	return e.InnerError
}

// Returns a new ResponseReadError
//
//	ResponseReadError: Occurs when there is an error from READING the response.
//	err: The inner error
func NewResponseReadError(err error) *ResponseReadError {
	return &ResponseReadError{
		InnerError: err,
	}
}

// ResponseParseError occurs when parsing the API response fails. E.g., invalid JSON.
type ResponseParseError struct {
	InnerError error
}

func (e *ResponseParseError) Error() string {
	return fmt.Sprintf("failed to parse API response: %v", e.InnerError)
}

func (e *ResponseParseError) ErrorCode() string {
	return "RESPONSE_PARSE_ERROR"
}

func (e *ResponseParseError) Unwrap() error {
	return e.InnerError
}

// Returns a new ResponseParseError
//
//	ResponseParseError: Occurs when there is an error PARSING the response.
//	err: The inner error
func NewResponseParseError(err error) *ResponseParseError {
	return &ResponseParseError{
		InnerError: err,
	}
}

// APIError occurs when API returns an explicit error message
type APIError struct {
	Message string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API error: %s", e.Message)
}

func (e *APIError) ErrorCode() string {
	return "API_ERROR"
}

// Returns a new APIError
//
//	APIError: Occurs when there is an explicit API error.
//	message: The explicit API error message.
func NewAPIError(message string) *APIError {
	return &APIError{
		Message: message,
	}
}

// DataNotFoundError occurs when expected data is missing from the API response
type DataNotFoundError struct {
	Key string // The key that was expected but not found
}

func (e *DataNotFoundError) Error() string {
	return fmt.Sprintf("could not find '%s' data in response.", e.Key)
}

func (e *DataNotFoundError) ErrorCode() string {
	return "DATA_NOT_FOUND_ERROR"
}

// Returns a new DataNotFoundError
//
//	DataNotFoundError: Occurs when the expected data is not found from the API.
//	key: The key that was expected but not found.
func NewDataNotFoundError(key string) *DataNotFoundError {
	return &DataNotFoundError{
		Key: key,
	}
}

// Compile time check to see if each error implements the ClientError interface
var (
	_ ClientError = (*HTTPRequestError)(nil)
	_ ClientError = (*HTTPStatusError)(nil)
	_ ClientError = (*ResponseReadError)(nil)
	_ ClientError = (*ResponseParseError)(nil)
	_ ClientError = (*APIError)(nil)
	_ ClientError = (*DataNotFoundError)(nil)
)
