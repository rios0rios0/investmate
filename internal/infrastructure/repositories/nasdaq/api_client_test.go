package nasdaq_test

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rios0rios0/investmate/internal/infrastructure/repositories/nasdaq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newAPIClient returns a client bound to an in-memory server that answers body to the request whose
// URI (path and query) is wantURI, and reports any other request as a test failure.
func newAPIClient(t *testing.T, wantURI string, body string) *nasdaq.APIClient {
	t.Helper()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, wantURI, r.URL.RequestURI())
		_, _ = io.WriteString(w, body)
	}))
	t.Cleanup(server.Close)

	return nasdaq.NewAPIClient(server.URL, server.Client())
}

// failingTransport is a stub round tripper whose every request fails before reaching a server.
type failingTransport struct{}

func (failingTransport) RoundTrip(*http.Request) (*http.Response, error) {
	return nil, errors.New("connection refused")
}

func TestAPIClient_FetchJSON(t *testing.T) {
	t.Parallel()

	t.Run("should decode the response when the request succeeds", func(t *testing.T) {
		t.Parallel()

		// given
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_ = json.NewEncoder(w).Encode(map[string]string{
				"method":    r.Method,
				"uri":       r.URL.RequestURI(),
				"userAgent": r.UserAgent(),
			})
		}))
		t.Cleanup(server.Close)
		client := nasdaq.NewAPIClient(server.URL+"/", server.Client())

		var got struct {
			Method    string `json:"method"`
			URI       string `json:"uri"`
			UserAgent string `json:"userAgent"`
		}

		// when
		err := client.FetchJSON(t.Context(), "/api/quote/SPY/dividends?assetclass=etf", &got)

		// then
		require.NoError(t, err)
		assert.Equal(t, http.MethodGet, got.Method)
		assert.Equal(t, "/api/quote/SPY/dividends?assetclass=etf", got.URI)
		assert.Contains(t, got.UserAgent, "Mozilla/5.0")
	})

	t.Run("should fall back to the default HTTP client when none is given", func(t *testing.T) {
		t.Parallel()

		// given
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			_, _ = io.WriteString(w, `{"ok": true}`)
		}))
		t.Cleanup(server.Close)
		client := nasdaq.NewAPIClient(server.URL, nil)

		var got struct {
			OK bool `json:"ok"`
		}

		// when
		err := client.FetchJSON(t.Context(), "/", &got)

		// then
		require.NoError(t, err)
		assert.True(t, got.OK)
	})

	errorCases := []struct {
		name    string
		client  func(t *testing.T) *nasdaq.APIClient
		wantErr string
	}{
		{
			name: "should return an error when the request cannot be created",
			client: func(*testing.T) *nasdaq.APIClient {
				return nasdaq.NewAPIClient("://invalid", http.DefaultClient)
			},
			wantErr: "failed to create request",
		},
		{
			name: "should return an error when the request cannot be sent",
			client: func(*testing.T) *nasdaq.APIClient {
				return nasdaq.NewAPIClient(nasdaq.DefaultBaseURL, &http.Client{Transport: failingTransport{}})
			},
			wantErr: "failed to fetch data",
		},
		{
			name: "should return an error when the response body is not JSON",
			client: func(t *testing.T) *nasdaq.APIClient {
				return newAPIClient(t, "/api/quote/SPY/dividends", "<html>blocked</html>")
			},
			wantErr: "failed to decode response",
		},
	}

	for _, testCase := range errorCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			// given
			client := testCase.client(t)

			var got map[string]any

			// when
			err := client.FetchJSON(t.Context(), "/api/quote/SPY/dividends", &got)

			// then
			require.ErrorContains(t, err, testCase.wantErr)
		})
	}
}
