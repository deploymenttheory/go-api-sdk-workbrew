package client

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/interfaces"
	"go.uber.org/zap"
	"resty.dev/v3"
)

// toInterfaceResponse converts a resty.Response to interfaces.Response.
// This internal function normalizes the response format for consistent handling.
//
// Parameters:
//   - resp: The resty response to convert (may be nil)
//
// Returns:
//   - *interfaces.Response: Normalized response with all metadata
func toInterfaceResponse(resp *resty.Response) *interfaces.Response {
	if resp == nil {
		return &interfaces.Response{
			Headers: make(http.Header),
		}
	}

	return &interfaces.Response{
		StatusCode: resp.StatusCode(),
		Status:     resp.Status(),
		Headers:    resp.Header(),
		Body:       []byte(resp.String()),
		Duration:   resp.Duration(),
		ReceivedAt: resp.ReceivedAt(),
		Size:       resp.Size(),
	}
}

// Response helper functions for working with interfaces.Response

// IsResponseSuccess returns true if the response status code indicates success (2xx).
//
// Parameters:
//   - resp: The response to check
//
// Returns:
//   - bool: True if status code is in the 200-299 range, false otherwise
func IsResponseSuccess(resp *interfaces.Response) bool {
	if resp == nil {
		return false
	}
	return resp.StatusCode >= 200 && resp.StatusCode < 300
}

// IsResponseError returns true if the response status code indicates an error (4xx or 5xx).
//
// Parameters:
//   - resp: The response to check
//
// Returns:
//   - bool: True if status code is 400 or higher, false otherwise
func IsResponseError(resp *interfaces.Response) bool {
	if resp == nil {
		return false
	}
	return resp.StatusCode >= 400
}

// GetResponseHeader retrieves a single header value from the response by key.
// Header lookup is case-insensitive following HTTP standards.
//
// Parameters:
//   - resp: The response containing headers
//   - key: The header name to retrieve (case-insensitive)
//
// Returns:
//   - string: The header value, or empty string if not found
//
// Example:
//
//	contentType := client.GetResponseHeader(resp, "Content-Type")
func GetResponseHeader(resp *interfaces.Response, key string) string {
	if resp == nil || resp.Headers == nil {
		return ""
	}
	return resp.Headers.Get(key)
}

// GetResponseHeaders returns all HTTP headers from the response.
//
// Parameters:
//   - resp: The response containing headers
//
// Returns:
//   - http.Header: All response headers, or empty map if response is nil
func GetResponseHeaders(resp *interfaces.Response) http.Header {
	if resp == nil {
		return make(http.Header)
	}
	return resp.Headers
}

// GetRateLimitHeaders extracts Workbrew API rate limiting headers from the response.
// These headers indicate the API quota limits and current usage.
//
// Parameters:
//   - resp: The response containing rate limit headers
//
// Returns:
//   - limit: Maximum number of requests allowed (X-Api-Quota-Limit)
//   - remaining: Number of requests remaining (X-Api-Quota-Remaining)
//   - reset: Unix timestamp when quota resets (X-Api-Quota-Reset)
//   - retryAfter: Seconds to wait before retrying (Retry-After)
//
// Example:
//
//	limit, remaining, reset, retryAfter := client.GetRateLimitHeaders(resp)
//	if remaining == "0" {
//	    log.Printf("Rate limit exceeded. Resets at: %s", reset)
//	    time.Sleep(time.Duration(retryAfter) * time.Second)
//	}
func GetRateLimitHeaders(resp *interfaces.Response) (limit, remaining, reset, retryAfter string) {
	if resp == nil {
		return
	}
	return resp.Headers.Get("X-Api-Quota-Limit"),
		resp.Headers.Get("X-Api-Quota-Remaining"),
		resp.Headers.Get("X-Api-Quota-Reset"),
		resp.Headers.Get("Retry-After")
}

// validateResponse validates the HTTP response before processing.
// This includes checking for empty responses and validating Content-Type for JSON endpoints.
//
// Parameters:
//   - resp: The resty response to validate
//   - method: The HTTP method used for logging
//   - path: The API endpoint path for logging
//
// Returns:
//   - error: Validation error if response is invalid, nil if valid
//
// Validation rules:
//   - Empty responses (204 No Content) are considered valid
//   - Non-error responses with content must have application/json Content-Type
//   - Error responses skip Content-Type validation (handled by error parser)
func (c *Client) validateResponse(resp *resty.Response, method, path string) error {
	// Handle empty responses (204 No Content, etc.)
	bodyLen := len(resp.String())
	if resp.Header().Get("Content-Length") == "0" || bodyLen == 0 {
		c.logger.Debug("Empty response received",
			zap.String("method", method),
			zap.String("path", path),
			zap.Int("status_code", resp.StatusCode()))
		return nil
	}

	// For non-error responses with content, validate Content-Type is JSON
	// Skip validation for:
	// - Error responses (handled by error parser)
	// - Endpoints that explicitly return non-JSON (download endpoints, etc.)
	if !resp.IsError() && bodyLen > 0 {
		contentType := resp.Header().Get("Content-Type")

		// Allow responses without Content-Type header (some endpoints don't set it)
		if contentType != "" && !strings.HasPrefix(contentType, "application/json") {
			c.logger.Warn("Unexpected Content-Type in response",
				zap.String("method", method),
				zap.String("path", path),
				zap.String("content_type", contentType),
				zap.String("expected", "application/json"))

			return fmt.Errorf("unexpected response Content-Type from %s %s: got %q, expected application/json",
				method, path, contentType)
		}
	}

	return nil
}
