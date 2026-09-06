package nasdaq_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/rios0rios0/investmate/internal/infrastructure/repositories/nasdaq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAPIPricesRepository_ListClosingPricesByETF(t *testing.T) {
	t.Parallel()

	currentYear := time.Now().Year()
	endpoint := fmt.Sprintf(
		"/api/quote/SPY/historical?assetclass=etf&fromdate=%d-01-01&todate=%d-12-31&limit=%d&offset=0",
		currentYear-nasdaq.YearsToFetch, currentYear, nasdaq.YearsToFetch*nasdaq.NumberOfDaysInYear,
	)

	t.Run("should average the closing prices per year when the API returns rows", func(t *testing.T) {
		t.Parallel()

		// given
		client := newAPIClient(t, endpoint, `{"data": {"tradesTable": {"rows": [
			{"close": "$100.00", "date": "01/02/2025"},
			{"close": "$110.00", "date": "01/03/2025"},
			{"close": "$90.00", "date": "12/30/2024"},
			{"close": "N/A", "date": "12/31/2024"}
		]}}}`)
		repository := nasdaq.NewAPIPricesRepositoryWithClient(client)

		// when
		prices, err := repository.ListClosingPricesByETF("SPY")

		// then
		require.NoError(t, err)
		assert.Len(t, prices, 2)
		assert.InDelta(t, 105.0, prices["2025"], 0.001)
		assert.InDelta(t, 90.0, prices["2024"], 0.001)
	})

	t.Run("should return an error when the API response cannot be decoded", func(t *testing.T) {
		t.Parallel()

		// given
		client := newAPIClient(t, endpoint, "<html>blocked</html>")
		repository := nasdaq.NewAPIPricesRepositoryWithClient(client)

		// when
		prices, err := repository.ListClosingPricesByETF("SPY")

		// then
		require.ErrorContains(t, err, "failed to decode response")
		assert.Nil(t, prices)
	})
}
