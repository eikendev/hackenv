// Package network provides HTTP client functionality.
package network

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// GetResponse performs an HTTP GET request and returns the response with proper error handling
func GetResponse(url string) (*http.Response, error) {
	resp, err := http.Get(url) //#nosec G107
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, fmt.Errorf("nil HTTP response from %s", url)
	}

	if resp.StatusCode != http.StatusOK {
		err = resp.Body.Close()
		if err != nil {
			log.Warnf("failed to close response body: %v", err)
		}

		return nil, fmt.Errorf("bad HTTP status code (%s)", resp.Status)
	}

	return resp, nil
}
