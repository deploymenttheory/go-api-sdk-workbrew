package client

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/interfaces"
	"go.uber.org/zap"
	"resty.dev/v3"
)

// Get executes a GET request
func (c *Client) Get(ctx context.Context, path string, queryParams map[string]string, headers map[string]string, result any) (*interfaces.Response, error) {
	req := c.client.R().
		SetContext(ctx).
		SetResult(result)

	for k, v := range queryParams {
		if v != "" {
			req.SetQueryParam(k, v)
		}
	}

	c.applyHeaders(req, headers)

	return c.executeRequest(req, "GET", path)
}

// Post executes a POST request with JSON body
func (c *Client) Post(ctx context.Context, path string, body any, headers map[string]string, result any) (*interfaces.Response, error) {
	req := c.client.R().
		SetContext(ctx).
		SetResult(result)

	if body != nil {
		req.SetBody(body)
	}

	c.applyHeaders(req, headers)

	return c.executeRequest(req, "POST", path)
}

// PostWithQuery executes a POST request with both body and query parameters
func (c *Client) PostWithQuery(ctx context.Context, path string, queryParams map[string]string, body any, headers map[string]string, result any) (*interfaces.Response, error) {
	req := c.client.R().
		SetContext(ctx).
		SetResult(result)

	for k, v := range queryParams {
		if v != "" {
			req.SetQueryParam(k, v)
		}
	}

	if body != nil {
		req.SetBody(body)
	}

	c.applyHeaders(req, headers)

	return c.executeRequest(req, "POST", path)
}

// Put executes a PUT request
func (c *Client) Put(ctx context.Context, path string, body any, headers map[string]string, result any) (*interfaces.Response, error) {
	req := c.client.R().
		SetContext(ctx).
		SetResult(result)

	if body != nil {
		req.SetBody(body)
	}

	c.applyHeaders(req, headers)

	return c.executeRequest(req, "PUT", path)
}

// Patch executes a PATCH request
func (c *Client) Patch(ctx context.Context, path string, body any, headers map[string]string, result any) (*interfaces.Response, error) {
	req := c.client.R().
		SetContext(ctx).
		SetResult(result)

	if body != nil {
		req.SetBody(body)
	}

	c.applyHeaders(req, headers)

	return c.executeRequest(req, "PATCH", path)
}

// Delete executes a DELETE request
func (c *Client) Delete(ctx context.Context, path string, queryParams map[string]string, headers map[string]string, result any) (*interfaces.Response, error) {
	req := c.client.R().
		SetContext(ctx).
		SetResult(result)

	for k, v := range queryParams {
		if v != "" {
			req.SetQueryParam(k, v)
		}
	}

	c.applyHeaders(req, headers)

	return c.executeRequest(req, "DELETE", path)
}

// DeleteWithBody executes a DELETE request with body (for bulk operations)
func (c *Client) DeleteWithBody(ctx context.Context, path string, body any, headers map[string]string, result any) (*interfaces.Response, error) {
	req := c.client.R().
		SetContext(ctx).
		SetResult(result)

	if body != nil {
		req.SetBody(body)
	}

	c.applyHeaders(req, headers)

	return c.executeRequest(req, "DELETE", path)
}

// GetCSV performs a GET request for CSV format and returns raw bytes
func (c *Client) GetCSV(ctx context.Context, path string, queryParams map[string]string, headers map[string]string) (*interfaces.Response, []byte, error) {
	var apiErr APIError
	req := c.client.R().
		SetContext(ctx).
		SetError(&apiErr)

	for k, v := range queryParams {
		if v != "" {
			req.SetQueryParam(k, v)
		}
	}

	c.applyHeaders(req, headers)

	c.logger.Debug("Executing CSV request",
		zap.String("method", "GET"),
		zap.String("path", path))

	resp, err := req.Get(path)
	ifaceResp := toInterfaceResponse(resp)
	if err != nil {
		c.logger.Error("CSV request failed",
			zap.String("path", path),
			zap.Error(err))
		return ifaceResp, nil, fmt.Errorf("CSV request failed: %w", err)
	}

	if resp.IsError() {
		return ifaceResp, nil, ParseErrorResponse(
			[]byte(resp.String()),
			resp.StatusCode(),
			resp.Status(),
			"GET",
			path,
			c.logger,
		)
	}

	body := []byte(resp.String())
	c.logger.Debug("CSV request completed successfully",
		zap.String("path", path),
		zap.Int("status_code", resp.StatusCode()),
		zap.Int("content_length", len(body)))

	return ifaceResp, body, nil
}

// executeRequest is a centralized request executor that handles error processing
// Returns response metadata and error. Response is always non-nil for accessing headers.
func (c *Client) executeRequest(req *resty.Request, method, path string) (*interfaces.Response, error) {
	c.logger.Debug("Executing API request",
		zap.String("method", method),
		zap.String("path", path))

	var resp *resty.Response
	var err error

	switch method {
	case "GET":
		resp, err = req.Get(path)
	case "POST":
		resp, err = req.Post(path)
	case "PUT":
		resp, err = req.Put(path)
	case "PATCH":
		resp, err = req.Patch(path)
	case "DELETE":
		resp, err = req.Delete(path)
	default:
		return toInterfaceResponse(nil), fmt.Errorf("unsupported HTTP method: %s", method)
	}

	// Convert to interface response (always return response metadata)
	ifaceResp := toInterfaceResponse(resp)

	if err != nil {
		c.logger.Error("Request failed",
			zap.String("method", method),
			zap.String("path", path),
			zap.Error(err))
		return ifaceResp, fmt.Errorf("request failed: %w", err)
	}

	// Validate response before processing
	if err := c.validateResponse(resp, method, path); err != nil {
		return ifaceResp, err
	}

	if resp.IsError() {
		return ifaceResp, ParseErrorResponse(
			[]byte(resp.String()),
			resp.StatusCode(),
			resp.Status(),
			method,
			path,
			c.logger,
		)
	}

	c.logger.Debug("Request completed successfully",
		zap.String("method", method),
		zap.String("path", path),
		zap.Int("status_code", resp.StatusCode()))

	return ifaceResp, nil
}
