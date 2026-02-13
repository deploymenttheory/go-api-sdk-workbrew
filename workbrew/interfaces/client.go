package interfaces

import (
	"context"
	"net/http"
	"time"

	"go.uber.org/zap"
)

// Response represents HTTP response metadata that can be returned alongside errors
// This allows callers to access response headers (rate limits, retry-after, etc.) even on errors
type Response struct {
	StatusCode int           // HTTP status code (e.g., 200, 404)
	Status     string        // HTTP status text (e.g., "200 OK")
	Headers    http.Header   // Response headers
	Body       []byte        // Raw response body
	Duration   time.Duration // Time taken for the request
	ReceivedAt time.Time     // When the response was received
	Size       int64         // Response body size in bytes
}

// HTTPClient interface that services will use
// This breaks import cycles by providing a contract without implementation
type HTTPClient interface {
	// Get executes a GET request and unmarshals the JSON response into the result parameter.
	// Returns response metadata and error. Response is non-nil even on error for accessing headers.
	Get(ctx context.Context, path string, queryParams map[string]string, headers map[string]string, result any) (*Response, error)
	
	// Post executes a POST request with a JSON body.
	// Returns response metadata and error. Response is non-nil even on error for accessing headers.
	Post(ctx context.Context, path string, body any, headers map[string]string, result any) (*Response, error)
	
	// PostWithQuery executes a POST request with both query parameters and a JSON body.
	// Returns response metadata and error. Response is non-nil even on error for accessing headers.
	PostWithQuery(ctx context.Context, path string, queryParams map[string]string, body any, headers map[string]string, result any) (*Response, error)
	
	// Put executes a PUT request with a JSON body.
	// Returns response metadata and error. Response is non-nil even on error for accessing headers.
	Put(ctx context.Context, path string, body any, headers map[string]string, result any) (*Response, error)
	
	// Patch executes a PATCH request with a JSON body.
	// Returns response metadata and error. Response is non-nil even on error for accessing headers.
	Patch(ctx context.Context, path string, body any, headers map[string]string, result any) (*Response, error)
	
	// Delete executes a DELETE request.
	// Returns response metadata and error. Response is non-nil even on error for accessing headers.
	Delete(ctx context.Context, path string, queryParams map[string]string, headers map[string]string, result any) (*Response, error)
	
	// DeleteWithBody executes a DELETE request with a JSON body (for bulk operations).
	// Returns response metadata and error. Response is non-nil even on error for accessing headers.
	DeleteWithBody(ctx context.Context, path string, body any, headers map[string]string, result any) (*Response, error)
	
	// GetCSV performs a GET request for CSV format and returns raw bytes.
	// Returns response metadata, CSV data, and error. Response is non-nil even on error for accessing headers.
	GetCSV(ctx context.Context, path string, queryParams map[string]string, headers map[string]string) (*Response, []byte, error)
	
	// GetLogger returns the configured zap logger instance.
	GetLogger() *zap.Logger
	
	// QueryBuilder returns a query builder instance for constructing URL query parameters.
	QueryBuilder() ServiceQueryBuilder
}

// ServiceQueryBuilder defines the query builder contract for services
type ServiceQueryBuilder interface {
	AddString(key, value string) QueryBuilder
	AddInt(key string, value int) QueryBuilder
	AddInt64(key string, value int64) QueryBuilder
	AddBool(key string, value bool) QueryBuilder
	AddTime(key string, value time.Time) QueryBuilder
	AddStringSlice(key string, values []string) QueryBuilder
	AddIntSlice(key string, values []int) QueryBuilder
	AddCustom(key, value string) QueryBuilder
	AddIfNotEmpty(key, value string) QueryBuilder
	AddIfTrue(condition bool, key, value string) QueryBuilder
	Merge(other map[string]string) QueryBuilder
	Remove(key string) QueryBuilder
	Has(key string) bool
	Get(key string) string
	Build() map[string]string
	BuildString() string
	Clear() QueryBuilder
	Count() int
	IsEmpty() bool
}

// QueryBuilder interface for method chaining
type QueryBuilder interface {
	ServiceQueryBuilder
}
