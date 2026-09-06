package nasdaq_test

import (
	"testing"

	"github.com/rios0rios0/investmate/internal/infrastructure/repositories/nasdaq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAPIDividendsRepository_ListDividendsByETF(t *testing.T) {
	t.Parallel()

	const endpoint = "/api/quote/SPY/dividends?assetclass=etf"

	t.Run("should sum the dividends per payment year when the API returns rows", func(t *testing.T) {
		t.Parallel()

		// given
		client := newAPIClient(t, endpoint, `{"data": {"dividends": {"rows": [
			{"amount": "$1.50", "paymentDate": "01/31/2025"},
			{"amount": "$2.25", "paymentDate": "04/30/2025"},
			{"amount": "$0.75", "paymentDate": "10/31/2024"},
			{"amount": "N/A", "paymentDate": "N/A"}
		]}}}`)
		repository := nasdaq.NewAPIDividendsRepositoryWithClient(client)

		// when
		dividends, err := repository.ListDividendsByETF("SPY")

		// then
		require.NoError(t, err)
		assert.Len(t, dividends, 2)
		assert.InDelta(t, 3.75, dividends["2025"], 0.001)
		assert.InDelta(t, 0.75, dividends["2024"], 0.001)
	})

	t.Run("should return an error when the API response cannot be decoded", func(t *testing.T) {
		t.Parallel()

		// given
		client := newAPIClient(t, endpoint, "<html>blocked</html>")
		repository := nasdaq.NewAPIDividendsRepositoryWithClient(client)

		// when
		dividends, err := repository.ListDividendsByETF("SPY")

		// then
		require.ErrorContains(t, err, "failed to decode response")
		assert.Nil(t, dividends)
	})
}
