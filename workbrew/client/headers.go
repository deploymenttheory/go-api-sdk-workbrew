package client

import "resty.dev/v3"

// applyHeaders applies HTTP headers to a request with proper precedence handling.
// This internal method ensures consistent header application across all request types.
//
// Header precedence (higher number wins):
//  1. Global headers (set via WithGlobalHeader/WithGlobalHeaders options)
//  2. Per-request headers (override global headers with the same key)
//
// Parameters:
//   - req: The resty request to apply headers to
//   - requestHeaders: Per-request headers that override global headers
//
// Example of header precedence:
//
//	// Global header set during client creation
//	client, _ := NewClient(apiKey, workspace, WithGlobalHeader("X-Custom", "global"))
//
//	// Per-request header overrides the global one
//	client.Get(ctx, "/endpoint", nil, map[string]string{"X-Custom": "request"}, &result)
//	// Final header: X-Custom: request (per-request value wins)
//
// Empty header values ("") are automatically filtered out to prevent sending empty headers.
func (t *Transport) applyHeaders(req *resty.Request, requestHeaders map[string]string) {
	// Apply global headers first
	for k, v := range t.globalHeaders {
		if v != "" {
			req.SetHeader(k, v)
		}
	}

	// Apply per-request headers (overrides global headers)
	for k, v := range requestHeaders {
		if v != "" {
			req.SetHeader(k, v)
		}
	}
}
