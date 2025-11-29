// Package network provides HTTP client functionality.
package network

import (
	"fmt"
	"log/slog"
	"net/http"
)

// GetResponse performs an HTTP GET request and returns the response with proper error handling
func GetResponse(url string) (*http.Response, error) {
	resp, err := http.Get(url) //#nosec G107
	if err != nil {
		return nil, fmt.Errorf("failed to perform HTTP GET %s: %w", url, err)
	}
	if resp == nil {
		return nil, fmt.Errorf("received nil HTTP response from %s", url)
	}

	if resp.StatusCode != http.StatusOK {
		err = resp.Body.Close()
		if err != nil {
			slog.Warn("Failed to close response body", "err", err, "url", url)
		}

		return nil, fmt.Errorf("received bad HTTP status code (%s)", resp.Status)
	}

	return resp, nil
}
