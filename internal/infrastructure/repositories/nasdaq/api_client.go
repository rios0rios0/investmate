package nasdaq

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const (
	// DefaultBaseURL is the public NASDAQ API origin the repositories talk to unless another one is given.
	DefaultBaseURL = "https://api.nasdaq.com"

	// userAgent imitates a desktop browser because the NASDAQ API answers requests without one with an empty body.
	userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 " +
		"(KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36 Edg/131.0.0.0"
)

// APIClient performs the JSON requests shared by every NASDAQ API repository: it builds the request
// with a context, sends the browser-like User-Agent the API requires and decodes the response body.
type APIClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewAPIClient returns an APIClient bound to the given API origin (scheme and host, with or without a
// trailing slash). A nil httpClient falls back to [http.DefaultClient].
func NewAPIClient(baseURL string, httpClient *http.Client) *APIClient {
	client := &APIClient{
		baseURL:    strings.TrimSuffix(baseURL, "/"),
		httpClient: httpClient,
	}

	if client.httpClient == nil {
		client.httpClient = http.DefaultClient
	}

	return client
}

// FetchJSON performs a GET request against the API endpoint (path plus query string) and decodes the
// JSON response body into out.
func (c *APIClient) FetchJSON(ctx context.Context, endpoint string, out any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+endpoint, http.NoBody)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", userAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to fetch data: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if err = json.NewDecoder(resp.Body).Decode(out); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	return nil
}
