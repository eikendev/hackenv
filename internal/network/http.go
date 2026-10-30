// Package network provides HTTP client functionality.
package network

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

// GetResponse performs an HTTP GET request and returns the response with proper error handling
func GetResponse(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to build HTTP GET %s: %w", url, err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform HTTP GET %s: %w", url, err)
	}
	if resp == nil {
		return nil, fmt.Errorf("received nil HTTP response from %s", url)
	}

	if resp.StatusCode != http.StatusOK {
		err = resp.Body.Close()
		if err != nil {
			slog.Warn("failed to close response body", "err", err, "url", url)
		}

		return nil, fmt.Errorf("received bad HTTP status code (%s)", resp.Status)
	}

	return resp, nil
}

// GetBytes performs an HTTP GET request and returns at most limit bytes of the body.
func GetBytes(ctx context.Context, url string, limit int64) ([]byte, error) {
	resp, err := GetResponse(ctx, url)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			slog.Warn("failed to close response body", "err", err, "url", url)
		}
	}()

	body, err := io.ReadAll(io.LimitReader(resp.Body, limit))
	if err != nil {
		return nil, fmt.Errorf("failed to read HTTP response from %s: %w", url, err)
	}

	return body, nil
}
