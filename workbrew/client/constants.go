package client

const (
	// DefaultBaseURL is the default base URL for the Workbrew API
	DefaultBaseURL = "https://console.workbrew.com"

	// DefaultAPIVersion is the default API version
	DefaultAPIVersion = "v0"

	// APIVersionHeader is the header name for the API version
	APIVersionHeader = "X-Workbrew-API-Version"

	// AuthorizationHeader is the header name for the authorization token
	AuthorizationHeader = "Authorization"

	// UserAgent is the user agent string for API requests
	UserAgent = "go-api-sdk-workbrew/1.0.0"

	// DefaultTimeout is the default HTTP client timeout in seconds
	DefaultTimeout = 120

	// MaxRetries is the maximum number of retries for failed requests
	MaxRetries = 3

	// RetryWaitTime is the wait time between retries in seconds
	RetryWaitTime = 2

	// RetryMaxWaitTime is the maximum wait time between retries in seconds
	RetryMaxWaitTime = 10
)

// Response format constants
const (
	FormatJSON = "json"
	FormatCSV  = "csv"
)
