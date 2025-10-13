package util

import (
	"context"
	"fmt"
	"io"
	"math"
	"net/http"
	"time"
)

// HTTPClient provides HTTP functionality with retry logic and timeout
type HTTPClient struct {
	client        *http.Client
	retryAttempts int
	timeout       time.Duration
	logger        *Logger
}

// NewHTTPClient creates a new HTTP client with configured timeout and retry settings
func NewHTTPClient(timeout time.Duration, retryAttempts int, logger *Logger) *HTTPClient {
	return &HTTPClient{
		client: &http.Client{
			Timeout: timeout,
		},
		retryAttempts: retryAttempts,
		timeout:       timeout,
		logger:        logger,
	}
}

// Get performs a GET request with retry logic and exponential backoff
func (h *HTTPClient) Get(ctx context.Context, url string, headers map[string]string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	return h.doWithRetry(req)
}

// Post performs a POST request with retry logic and exponential backoff
func (h *HTTPClient) Post(ctx context.Context, url string, body io.Reader, headers map[string]string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	return h.doWithRetry(req)
}

// doWithRetry executes an HTTP request with exponential backoff retry logic
func (h *HTTPClient) doWithRetry(req *http.Request) (*http.Response, error) {
	var resp *http.Response
	var err error

	for attempt := 0; attempt <= h.retryAttempts; attempt++ {
		resp, err = h.client.Do(req)

		// Request succeeded
		if err == nil && resp.StatusCode < 500 {
			// Check for rate limiting (429)
			if resp.StatusCode == http.StatusTooManyRequests {
				if attempt < h.retryAttempts {
					h.logger.Warn("Rate limited, retrying",
						"attempt", attempt+1,
						"url", req.URL.String())
					_ = resp.Body.Close()
					h.backoff(attempt)
					continue
				}
				return resp, fmt.Errorf("rate limit exceeded after %d attempts", attempt+1)
			}
			return resp, nil
		} // Log the error
		if err != nil {
			h.logger.Warn("HTTP request failed",
				"attempt", attempt+1,
				"url", req.URL.String(),
				"error", err.Error())
		} else {
			h.logger.Warn("HTTP request returned server error",
				"attempt", attempt+1,
				"url", req.URL.String(),
				"status", resp.StatusCode)
			_ = resp.Body.Close()
		} // Don't retry on last attempt
		if attempt < h.retryAttempts {
			h.backoff(attempt)
		}
	}

	if err != nil {
		return nil, fmt.Errorf("request failed after %d attempts: %w", h.retryAttempts+1, err)
	}

	return nil, fmt.Errorf("request failed after %d attempts with status %d", h.retryAttempts+1, resp.StatusCode)
}

// backoff implements exponential backoff with jitter
func (h *HTTPClient) backoff(attempt int) {
	// Exponential backoff: 1s, 2s, 4s, 8s, etc.
	backoffDuration := time.Duration(math.Pow(2, float64(attempt))) * time.Second

	// Cap at 30 seconds
	if backoffDuration > 30*time.Second {
		backoffDuration = 30 * time.Second
	}

	h.logger.Debug("Backing off before retry",
		"attempt", attempt+1,
		"duration", backoffDuration.String())

	time.Sleep(backoffDuration)
}

// ReadResponseBody reads and closes the response body
func ReadResponseBody(resp *http.Response) ([]byte, error) {
	defer func() {
		_ = resp.Body.Close()
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return body, nil
}

// CheckResponse checks if an HTTP response indicates an error
func CheckResponse(resp *http.Response) error {
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}

	body, err := ReadResponseBody(resp)
	if err != nil {
		return fmt.Errorf("HTTP %d: failed to read error response", resp.StatusCode)
	}

	return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
}
