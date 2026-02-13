package client

import (
	"encoding/json"
	"fmt"
	"strings"

	"go.uber.org/zap"
)

// APIError represents an error response from the Workbrew API.
// It contains both the high-level error message and detailed validation errors.
//
// The error structure matches the Workbrew API swagger specification:
//
//	{
//	  "message": "string",
//	  "errors": ["string"]
//	}
//
// Example error response:
//
//	{
//	  "message": "Validation failed",
//	  "errors": [
//	    "Field 'name' is required",
//	    "Field 'content' must not be empty"
//	  ]
//	}
//
// Use the helper functions (IsUnauthorized, IsNotFound, etc.) to check for specific error types.
type APIError struct {
	Message string   `json:"message"`          // Main error message
	Errors  []string `json:"errors,omitempty"` // Array of detailed error messages

	// HTTP response details
	StatusCode int    // HTTP status code
	Status     string // HTTP status text
	Endpoint   string // API endpoint that returned the error
	Method     string // HTTP method used
}

// Error implements the error interface.
// It formats the error message with context including HTTP method, endpoint, and status code.
//
// Returns:
//   - string: Formatted error message with full context
//
// Example output:
//
//	"Workbrew API error (422 Unprocessable Entity) at POST /api/v1/brewfiles: Validation failed - [Field 'name' is required]"
func (e *APIError) Error() string {
	if len(e.Errors) > 0 {
		return fmt.Sprintf("Workbrew API error (%d %s) at %s %s: %s - %v",
			e.StatusCode, e.Status, e.Method, e.Endpoint, e.Message, e.Errors)
	}
	return fmt.Sprintf("Workbrew API error (%d %s) at %s %s: %s",
		e.StatusCode, e.Status, e.Method, e.Endpoint, e.Message)
}

// ParseErrorResponse parses an HTTP error response body into an APIError.
// It attempts to parse the response as JSON and falls back to using the raw body as the message.
//
// Parameters:
//   - body: Raw response body bytes
//   - statusCode: HTTP status code
//   - status: HTTP status text (e.g., "404 Not Found")
//   - method: HTTP method used
//   - endpoint: API endpoint path
//   - logger: Logger for recording error details
//
// Returns:
//   - error: APIError containing parsed error information
//
// Example usage:
//
//	if resp.IsError() {
//	    return ParseErrorResponse(resp.Body(), resp.StatusCode(), resp.Status(), "GET", "/api/brewfiles", logger)
//	}
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

// getDefaultErrorMessage returns a descriptive default error message for HTTP status codes.
// Used when the API response doesn't contain a message or isn't valid JSON.
//
// Parameters:
//   - statusCode: HTTP status code
//
// Returns:
//   - string: Human-readable error message with actionable guidance
func getDefaultErrorMessage(statusCode int) string {
	switch statusCode {
	case StatusBadRequest:
		return "Bad request - the request is invalid or malformed"
	case StatusUnauthorized:
		return "Authentication required or invalid API key. Verify that you have provided your correct API key."
	case StatusForbidden:
		return "Access forbidden - you are not allowed to perform this operation. May require plan upgrade."
	case StatusNotFound:
		return "Resource not found"
	case StatusConflict:
		return "Resource already exists"
	case StatusUnprocessableEntity:
		return "Validation error - the request contains invalid parameters"
	case StatusFailedDependency:
		return "The request depended on another request that failed"
	case StatusTooManyRequests:
		return "Rate limit exceeded. Too many requests have been made in a given amount of time. Please retry after some time."
	case StatusInternalServerError:
		return "Internal server error"
	case StatusBadGateway:
		return "Bad gateway"
	case StatusServiceUnavailable:
		return "Service temporarily unavailable. Retry might work."
	case StatusGatewayTimeout:
		return "The operation took too long to complete. Request timeout."
	default:
		return "Unknown error"
	}
}

// Error type check helpers
//
// The following functions provide convenient type checking for specific API error conditions.
// They help build resilient error handling logic by checking for specific HTTP status codes.

// IsBadRequest checks if the error is a bad request error (400).
// This typically indicates invalid request parameters or malformed JSON.
//
// Parameters:
//   - err: The error to check
//
// Returns:
//   - bool: True if the error is a 400 Bad Request, false otherwise
//
// Example:
//
//	if client.IsBadRequest(err) {
//	    log.Println("Invalid request parameters:", err)
//	}
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

// IsConflict checks if the error is a conflict error (409) - resource already exists
func IsConflict(err error) bool {
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.StatusCode == StatusConflict
	}
	return false
}

// IsRateLimited checks if the error is a rate limit error (429).
// This indicates too many requests have been made within the rate limit window.
//
// Parameters:
//   - err: The error to check
//
// Returns:
//   - bool: True if the error is a 429 Too Many Requests, false otherwise
//
// Example with retry logic:
//
//	if client.IsRateLimited(err) {
//	    _, _, reset, retryAfter := client.GetRateLimitHeaders(resp)
//	    log.Printf("Rate limited. Retry after: %s seconds", retryAfter)
//	    time.Sleep(time.Duration(retryAfterInt) * time.Second)
//	    // Retry the request
//	}
func IsRateLimited(err error) bool {
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.StatusCode == StatusTooManyRequests
	}
	return false
}

// IsTransient checks if the error is transient and safe to retry.
// Returns true for 503 (Service Unavailable) and 504 (Gateway Timeout).
// These errors are typically temporary and the request may succeed on retry.
//
// Parameters:
//   - err: The error to check
//
// Returns:
//   - bool: True if the error is retryable (503 or 504), false otherwise
//
// Example with exponential backoff:
//
//	if client.IsTransient(err) {
//	    for attempt := 0; attempt < 3; attempt++ {
//	        time.Sleep(time.Duration(math.Pow(2, float64(attempt))) * time.Second)
//	        result, _, err := client.Get(ctx, path, params, headers, &data)
//	        if err == nil || !client.IsTransient(err) {
//	            break
//	        }
//	    }
//	}
func IsTransient(err error) bool {
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.StatusCode == StatusServiceUnavailable ||
			apiErr.StatusCode == StatusGatewayTimeout
	}
	return false
}

// IsDeadlineExceeded checks if the operation took too long to complete (504)
func IsDeadlineExceeded(err error) bool {
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.StatusCode == StatusGatewayTimeout
	}
	return false
}
