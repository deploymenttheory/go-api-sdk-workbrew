package client

import (
	"encoding/json"
	"fmt"
	"strings"

	"go.uber.org/zap"
)

// HTTP Status Codes from Workbrew API Swagger Specification
const (
	// Success status codes
	StatusOK      = 200 // Successful GET requests
	StatusCreated = 201 // Successful POST requests (brew_commands, brewfiles)

	// Client error status codes
	StatusBadRequest          = 400 // Bad request (general client error)
	StatusUnauthorized        = 401 // Invalid or missing authentication token
	StatusForbidden           = 403 // Free tier plan restrictions, insufficient permissions
	StatusNotFound            = 404 // Resource not found
	StatusUnprocessableEntity = 422 // Validation errors, invalid request parameters

	// Server error status codes
	StatusInternalServerError = 500 // Server-side error
	StatusBadGateway          = 502 // Gateway error
	StatusServiceUnavailable  = 503 // Service temporarily unavailable
)

// APIError represents an error response from the Workbrew API
// Matches the error schema from swagger:
//
//	{
//	  "message": "string",
//	  "errors": ["string"]
//	}
type APIError struct {
	Message string   `json:"message"`           // Main error message
	Errors  []string `json:"errors,omitempty"` // Array of detailed error messages

	// HTTP response details
	StatusCode int    // HTTP status code
	Status     string // HTTP status text
	Endpoint   string // API endpoint that returned the error
	Method     string // HTTP method used
}

// Error implements the error interface
func (e *APIError) Error() string {
	if len(e.Errors) > 0 {
		return fmt.Sprintf("Workbrew API error (%d %s) at %s %s: %s - %v",
			e.StatusCode, e.Status, e.Method, e.Endpoint, e.Message, e.Errors)
	}
	return fmt.Sprintf("Workbrew API error (%d %s) at %s %s: %s",
		e.StatusCode, e.Status, e.Method, e.Endpoint, e.Message)
}

// ParseErrorResponse parses an error response from the API
func ParseErrorResponse(body []byte, statusCode int, status, method, endpoint string, logger *zap.Logger) error {
	apiError := &APIError{
		StatusCode: statusCode,
		Status:     status,
		Endpoint:   endpoint,
		Method:     method,
	}

	// Try to parse as JSON error response
	if err := json.Unmarshal(body, apiError); err != nil {
		// If JSON parsing fails, use the raw body as message
		apiError.Message = string(body)
		logger.Debug("Failed to parse error response as JSON, using raw body",
			zap.Error(err),
			zap.String("body", string(body)))
	}

	// If no message was parsed, set a default message based on status code
	if apiError.Message == "" {
		apiError.Message = getDefaultErrorMessage(statusCode)
	}

	logger.Error("API error response",
		zap.Int("status_code", statusCode),
		zap.String("status", status),
		zap.String("method", method),
		zap.String("endpoint", endpoint),
		zap.String("message", apiError.Message),
		zap.Strings("errors", apiError.Errors))

	return apiError
}

// getDefaultErrorMessage returns a default error message based on status code
func getDefaultErrorMessage(statusCode int) string {
	switch statusCode {
	case StatusBadRequest:
		return "Bad request"
	case StatusUnauthorized:
		return "Authentication required or invalid API key"
	case StatusForbidden:
		return "Access forbidden - may require plan upgrade"
	case StatusNotFound:
		return "Resource not found"
	case StatusUnprocessableEntity:
		return "Validation error"
	case StatusInternalServerError:
		return "Internal server error"
	case StatusBadGateway:
		return "Bad gateway"
	case StatusServiceUnavailable:
		return "Service temporarily unavailable"
	default:
		return "Unknown error"
	}
}

// Error type check helpers

// IsBadRequest checks if the error is a bad request error (400)
func IsBadRequest(err error) bool {
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.StatusCode == StatusBadRequest
	}
	return false
}

// IsUnauthorized checks if the error is an authentication error (401)
func IsUnauthorized(err error) bool {
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.StatusCode == StatusUnauthorized
	}
	return false
}

// IsForbidden checks if the error is a forbidden error (403)
// This typically indicates free tier plan restrictions per swagger spec
func IsForbidden(err error) bool {
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.StatusCode == StatusForbidden
	}
	return false
}

// IsNotFound checks if the error is a not found error (404)
func IsNotFound(err error) bool {
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.StatusCode == StatusNotFound
	}
	return false
}

// IsValidationError checks if the error is a validation/unprocessable entity error (422)
// Per swagger: "Arguments cannot include `&&`", "Brewfile has an invalid line", etc.
func IsValidationError(err error) bool {
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.StatusCode == StatusUnprocessableEntity
	}
	return false
}

// IsServerError checks if the error is a server error (5xx)
func IsServerError(err error) bool {
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.StatusCode >= 500 && apiErr.StatusCode < 600
	}
	return false
}

// IsFreeTierError checks if the error is specifically a free tier restriction error (403)
// Per swagger spec, these errors have messages about plan upgrades
func IsFreeTierError(err error) bool {
	if apiErr, ok := err.(*APIError); ok {
		if apiErr.StatusCode == StatusForbidden {
			// Check if error message contains free tier related text
			for _, errMsg := range apiErr.Errors {
				if strings.Contains(strings.ToLower(errMsg), "free subscription") ||
					strings.Contains(strings.ToLower(errMsg), "upgrade your plan") {
					return true
				}
			}
			if strings.Contains(strings.ToLower(apiErr.Message), "free subscription") ||
				strings.Contains(strings.ToLower(apiErr.Message), "upgrade your plan") {
				return true
			}
		}
	}
	return false
}
