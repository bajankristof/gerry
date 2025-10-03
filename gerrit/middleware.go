package gerrit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"resty.dev/v3"
)

// autoErrorMiddleware automatically returns an error for non-2xx responses
func autoErrorMiddleware(c *resty.Client, r *resty.Response) error {
	if r.StatusCode() < 200 || r.StatusCode() >= 300 {
		return fmt.Errorf("Gerrit API error: %d %s", r.StatusCode(), r.Status())
	}

	return nil
}

// autoParseMiddleware automatically parses the response body into the Result field of the request
func autoParseMiddleware(c *resty.Client, r *resty.Response) error {
	if r.Err != nil || r.IsError() {
		return nil
	}

	if r.Request.Result == nil {
		return nil
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if len(b) == 0 {
		return nil
	}

	// Strip Gerrit's XSSI protection prefix
	b = bytes.TrimPrefix(b, []byte(")]}'\n"))

	if err := json.Unmarshal(b, r.Request.Result); err != nil {
		return fmt.Errorf("failed to parse JSON response: %w", err)
	}

	return nil
}
