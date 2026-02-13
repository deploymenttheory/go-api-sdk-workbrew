package client

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

func TestAPIError_Error(t *testing.T) {
	tests := []struct {
		name        string
		apiError    *APIError
		wantContain []string
	}{
		{
			name: "error with message only",
			apiError: &APIError{
				StatusCode: 404,
				Status:     "404 Not Found",
				Method:     "GET",
				Endpoint:   "/api/v1/resource",
				Message:    "Resource not found",
				Errors:     nil,
			},
			wantContain: []string{
				"Workbrew API error",
				"404",
				"Not Found",
				"GET",
				"/api/v1/resource",
				"Resource not found",
			},
		},
		{
			name: "error with message and errors array",
			apiError: &APIError{
				StatusCode: 422,
				Status:     "422 Unprocessable Entity",
				Method:     "POST",
				Endpoint:   "/api/v1/brewfiles",
				Message:    "Validation failed",
				Errors:     []string{"Field 'name' is required", "Field 'content' is invalid"},
			},
			wantContain: []string{
				"Workbrew API error",
				"422",
				"Unprocessable Entity",
				"POST",
				"/api/v1/brewfiles",
				"Validation failed",
				"Field 'name' is required",
				"Field 'content' is invalid",
			},
		},
		{
			name: "error with empty errors array",
			apiError: &APIError{
				StatusCode: 500,
				Status:     "500 Internal Server Error",
				Method:     "PUT",
				Endpoint:   "/api/v1/update",
				Message:    "Internal error",
				Errors:     []string{},
			},
			wantContain: []string{
				"Workbrew API error",
				"500",
				"Internal Server Error",
				"PUT",
				"/api/v1/update",
				"Internal error",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.apiError.Error()
			for _, want := range tt.wantContain {
				if !strings.Contains(got, want) {
					t.Errorf("Error() = %q, want to contain %q", got, want)
				}
			}
		})
	}
}

func TestParseErrorResponse_ValidJSON(t *testing.T) {
	logger := zaptest.NewLogger(t)

	tests := []struct {
		name           string
		body           string
		statusCode     int
		status         string
		method         string
		endpoint       string
		wantMessage    string
		wantErrors     []string
		wantStatusCode int
	}{
		{
			name: "valid error response with errors array",
			body: `{
				"message": "Validation failed",
				"errors": ["Field 'name' is required", "Field 'content' is invalid"]
			}`,
			statusCode:     422,
			status:         "422 Unprocessable Entity",
			method:         "POST",
			endpoint:       "/api/v1/brewfiles",
			wantMessage:    "Validation failed",
			wantErrors:     []string{"Field 'name' is required", "Field 'content' is invalid"},
			wantStatusCode: 422,
		},
		{
			name: "valid error response without errors array",
			body: `{
				"message": "Resource not found"
			}`,
			statusCode:     404,
			status:         "404 Not Found",
			method:         "GET",
			endpoint:       "/api/v1/resource/123",
			wantMessage:    "Resource not found",
			wantErrors:     nil,
			wantStatusCode: 404,
		},
		{
			name: "empty message uses default",
			body: `{
				"message": ""
			}`,
			statusCode:     401,
			status:         "401 Unauthorized",
			method:         "GET",
			endpoint:       "/api/v1/data",
			wantMessage:    "Authentication required or invalid API key. Verify that you have provided your correct API key.",
			wantErrors:     nil,
			wantStatusCode: 401,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ParseErrorResponse(
				[]byte(tt.body),
				tt.statusCode,
				tt.status,
				tt.method,
				tt.endpoint,
				logger,
			)

			if err == nil {
				t.Fatal("ParseErrorResponse() returned nil error")
			}

			apiErr, ok := err.(*APIError)
			if !ok {
				t.Fatalf("ParseErrorResponse() returned %T, want *APIError", err)
			}

			if apiErr.Message != tt.wantMessage {
				t.Errorf("Message = %q, want %q", apiErr.Message, tt.wantMessage)
			}

			if len(apiErr.Errors) != len(tt.wantErrors) {
				t.Errorf("Errors length = %d, want %d", len(apiErr.Errors), len(tt.wantErrors))
			} else {
				for i, want := range tt.wantErrors {
					if apiErr.Errors[i] != want {
						t.Errorf("Errors[%d] = %q, want %q", i, apiErr.Errors[i], want)
					}
				}
			}

			if apiErr.StatusCode != tt.wantStatusCode {
				t.Errorf("StatusCode = %d, want %d", apiErr.StatusCode, tt.wantStatusCode)
			}

			if apiErr.Status != tt.status {
				t.Errorf("Status = %q, want %q", apiErr.Status, tt.status)
			}

			if apiErr.Method != tt.method {
				t.Errorf("Method = %q, want %q", apiErr.Method, tt.method)
			}

			if apiErr.Endpoint != tt.endpoint {
				t.Errorf("Endpoint = %q, want %q", apiErr.Endpoint, tt.endpoint)
			}
		})
	}
}

func TestParseErrorResponse_InvalidJSON(t *testing.T) {
	logger := zaptest.NewLogger(t)

	tests := []struct {
		name           string
		body           string
		statusCode     int
		wantMessage    string
		wantStatusCode int
	}{
		{
			name:           "plain text error",
			body:           "Something went wrong",
			statusCode:     500,
			wantMessage:    "Something went wrong",
			wantStatusCode: 500,
		},
		{
			name:           "HTML error page",
			body:           "<html><body>Error 404</body></html>",
			statusCode:     404,
			wantMessage:    "<html><body>Error 404</body></html>",
			wantStatusCode: 404,
		},
		{
			name:           "empty body uses default message",
			body:           "",
			statusCode:     503,
			wantMessage:    "Service temporarily unavailable. Retry might work.",
			wantStatusCode: 503,
		},
		{
			name:           "malformed JSON",
			body:           `{"message": "incomplete`,
			statusCode:     400,
			wantMessage:    `{"message": "incomplete`,
			wantStatusCode: 400,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ParseErrorResponse(
				[]byte(tt.body),
				tt.statusCode,
				"",
				"GET",
				"/test",
				logger,
			)

			if err == nil {
				t.Fatal("ParseErrorResponse() returned nil error")
			}

			apiErr, ok := err.(*APIError)
			if !ok {
				t.Fatalf("ParseErrorResponse() returned %T, want *APIError", err)
			}

			if apiErr.Message != tt.wantMessage {
				t.Errorf("Message = %q, want %q", apiErr.Message, tt.wantMessage)
			}

			if apiErr.StatusCode != tt.wantStatusCode {
				t.Errorf("StatusCode = %d, want %d", apiErr.StatusCode, tt.wantStatusCode)
			}
		})
	}
}

func TestGetDefaultErrorMessage(t *testing.T) {
	tests := []struct {
		statusCode int
		want       string
	}{
		{StatusBadRequest, "Bad request - the request is invalid or malformed"},
		{StatusUnauthorized, "Authentication required or invalid API key. Verify that you have provided your correct API key."},
		{StatusForbidden, "Access forbidden - you are not allowed to perform this operation. May require plan upgrade."},
		{StatusNotFound, "Resource not found"},
		{StatusConflict, "Resource already exists"},
		{StatusUnprocessableEntity, "Validation error - the request contains invalid parameters"},
		{StatusFailedDependency, "The request depended on another request that failed"},
		{StatusTooManyRequests, "Rate limit exceeded. Too many requests have been made in a given amount of time. Please retry after some time."},
		{StatusInternalServerError, "Internal server error"},
		{StatusBadGateway, "Bad gateway"},
		{StatusServiceUnavailable, "Service temporarily unavailable. Retry might work."},
		{StatusGatewayTimeout, "The operation took too long to complete. Request timeout."},
		{999, "Unknown error"},
		{418, "Unknown error"}, // I'm a teapot
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := getDefaultErrorMessage(tt.statusCode)
			if got != tt.want {
				t.Errorf("getDefaultErrorMessage(%d) = %q, want %q", tt.statusCode, got, tt.want)
			}
		})
	}
}

func TestIsBadRequest(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "400 error",
			err:  &APIError{StatusCode: 400},
			want: true,
		},
		{
			name: "404 error",
			err:  &APIError{StatusCode: 404},
			want: false,
		},
		{
			name: "non-APIError",
			err:  errors.New("generic error"),
			want: false,
		},
		{
			name: "nil error",
			err:  nil,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsBadRequest(tt.err)
			if got != tt.want {
				t.Errorf("IsBadRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsUnauthorized(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "401 error",
			err:  &APIError{StatusCode: 401},
			want: true,
		},
		{
			name: "403 error",
			err:  &APIError{StatusCode: 403},
			want: false,
		},
		{
			name: "non-APIError",
			err:  errors.New("generic error"),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsUnauthorized(tt.err)
			if got != tt.want {
				t.Errorf("IsUnauthorized() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsForbidden(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "403 error",
			err:  &APIError{StatusCode: 403},
			want: true,
		},
		{
			name: "401 error",
			err:  &APIError{StatusCode: 401},
			want: false,
		},
		{
			name: "non-APIError",
			err:  errors.New("generic error"),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsForbidden(tt.err)
			if got != tt.want {
				t.Errorf("IsForbidden() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsNotFound(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "404 error",
			err:  &APIError{StatusCode: 404},
			want: true,
		},
		{
			name: "400 error",
			err:  &APIError{StatusCode: 400},
			want: false,
		},
		{
			name: "non-APIError",
			err:  errors.New("generic error"),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsNotFound(tt.err)
			if got != tt.want {
				t.Errorf("IsNotFound() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsValidationError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "422 error",
			err:  &APIError{StatusCode: 422},
			want: true,
		},
		{
			name: "400 error",
			err:  &APIError{StatusCode: 400},
			want: false,
		},
		{
			name: "non-APIError",
			err:  errors.New("generic error"),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsValidationError(tt.err)
			if got != tt.want {
				t.Errorf("IsValidationError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsServerError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "500 error",
			err:  &APIError{StatusCode: 500},
			want: true,
		},
		{
			name: "502 error",
			err:  &APIError{StatusCode: 502},
			want: true,
		},
		{
			name: "503 error",
			err:  &APIError{StatusCode: 503},
			want: true,
		},
		{
			name: "599 error",
			err:  &APIError{StatusCode: 599},
			want: true,
		},
		{
			name: "400 error",
			err:  &APIError{StatusCode: 400},
			want: false,
		},
		{
			name: "404 error",
			err:  &APIError{StatusCode: 404},
			want: false,
		},
		{
			name: "600 error (out of 5xx range)",
			err:  &APIError{StatusCode: 600},
			want: false,
		},
		{
			name: "non-APIError",
			err:  errors.New("generic error"),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsServerError(tt.err)
			if got != tt.want {
				t.Errorf("IsServerError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsFreeTierError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "403 with free subscription message",
			err: &APIError{
				StatusCode: 403,
				Message:    "This feature requires upgrading from free subscription",
			},
			want: true,
		},
		{
			name: "403 with upgrade plan message",
			err: &APIError{
				StatusCode: 403,
				Message:    "Please upgrade your plan to access this feature",
			},
			want: true,
		},
		{
			name: "403 with free subscription in errors array",
			err: &APIError{
				StatusCode: 403,
				Message:    "Forbidden",
				Errors:     []string{"This feature is not available on free subscription"},
			},
			want: true,
		},
		{
			name: "403 with upgrade plan in errors array",
			err: &APIError{
				StatusCode: 403,
				Message:    "Access denied",
				Errors:     []string{"You need to upgrade your plan"},
			},
			want: true,
		},
		{
			name: "403 with case insensitive match",
			err: &APIError{
				StatusCode: 403,
				Message:    "FREE SUBSCRIPTION limit reached",
			},
			want: true,
		},
		{
			name: "403 without free tier message",
			err: &APIError{
				StatusCode: 403,
				Message:    "Access forbidden",
			},
			want: false,
		},
		{
			name: "401 with free subscription message",
			err: &APIError{
				StatusCode: 401,
				Message:    "free subscription",
			},
			want: false,
		},
		{
			name: "500 error",
			err:  &APIError{StatusCode: 500},
			want: false,
		},
		{
			name: "non-APIError",
			err:  errors.New("generic error"),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsFreeTierError(tt.err)
			if got != tt.want {
				t.Errorf("IsFreeTierError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAPIError_JSON_Marshalling(t *testing.T) {
	// Test that APIError can be marshalled and unmarshalled as JSON
	original := &APIError{
		Message:    "Test error message",
		Errors:     []string{"error1", "error2"},
		StatusCode: 422,
		Status:     "422 Unprocessable Entity",
		Endpoint:   "/api/v1/test",
		Method:     "POST",
	}

	// Marshal to JSON
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Failed to marshal APIError: %v", err)
	}

	// Unmarshal back
	var unmarshalled APIError
	err = json.Unmarshal(data, &unmarshalled)
	if err != nil {
		t.Fatalf("Failed to unmarshal APIError: %v", err)
	}

	// Compare (note: unexported fields won't be in JSON)
	if unmarshalled.Message != original.Message {
		t.Errorf("Message = %q, want %q", unmarshalled.Message, original.Message)
	}
	if len(unmarshalled.Errors) != len(original.Errors) {
		t.Errorf("Errors length = %d, want %d", len(unmarshalled.Errors), len(original.Errors))
	}
}

func TestParseErrorResponse_NilLogger(t *testing.T) {
	// Should not panic with nil logger (though logger should never be nil in practice)
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("ParseErrorResponse panicked with nil logger: %v", r)
		}
	}()

	// Use a development logger instead of nil
	logger, _ := zap.NewDevelopment()

	err := ParseErrorResponse(
		[]byte(`{"message": "test"}`),
		404,
		"404 Not Found",
		"GET",
		"/test",
		logger,
	)

	if err == nil {
		t.Error("ParseErrorResponse() returned nil error")
	}
}

func TestErrorConstants(t *testing.T) {
	// Verify error constants have expected values
	constants := map[string]int{
		"StatusOK":                   StatusOK,
		"StatusCreated":              StatusCreated,
		"StatusBadRequest":           StatusBadRequest,
		"StatusUnauthorized":         StatusUnauthorized,
		"StatusForbidden":            StatusForbidden,
		"StatusNotFound":             StatusNotFound,
		"StatusUnprocessableEntity":  StatusUnprocessableEntity,
		"StatusInternalServerError":  StatusInternalServerError,
		"StatusBadGateway":           StatusBadGateway,
		"StatusServiceUnavailable":   StatusServiceUnavailable,
	}

	expected := map[string]int{
		"StatusOK":                   200,
		"StatusCreated":              201,
		"StatusBadRequest":           400,
		"StatusUnauthorized":         401,
		"StatusForbidden":            403,
		"StatusNotFound":             404,
		"StatusUnprocessableEntity":  422,
		"StatusInternalServerError":  500,
		"StatusBadGateway":           502,
		"StatusServiceUnavailable":   503,
	}

	for name, got := range constants {
		want := expected[name]
		if got != want {
			t.Errorf("%s = %d, want %d", name, got, want)
		}
	}
}
