package main

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type stubDividendsRepository struct {
	data map[string]float64
	err  error
}

func (s *stubDividendsRepository) ListDividendsByETF(_ string) (map[string]float64, error) {
	return s.data, s.err
}

type stubPricesRepository struct {
	data map[string]float64
	err  error
}

func (s *stubPricesRepository) ListClosingPricesByETF(_ string) (map[string]float64, error) {
	return s.data, s.err
}

func TestMain_ProcessETF(t *testing.T) {
	t.Parallel()

	t.Run("should populate ETF with data when repositories return valid results", func(t *testing.T) {
		t.Parallel()

		// given
		dividendsRepo := &stubDividendsRepository{
			data: map[string]float64{"2025": 5.50, "2024": 4.80},
		}
		pricesRepo := &stubPricesRepository{
			data: map[string]float64{"2025": 450.00, "2024": 420.00},
		}

		// when
		etf := processETF("SPY", dividendsRepo, pricesRepo)

		// then
		assert.Equal(t, "SPY", etf.Name)
		assert.NotEmpty(t, etf.AmountDividendsPerYear)
		assert.NotEmpty(t, etf.AverageClosingPricePerYear)
		assert.InDelta(t, 5.50, etf.AmountDividendsPerYear["2025"], 0.001)
		assert.InDelta(t, 450.00, etf.AverageClosingPricePerYear["2025"], 0.001)
	})

	t.Run("should return ETF with empty maps when repositories return errors", func(t *testing.T) {
		t.Parallel()

		// given
		dividendsRepo := &stubDividendsRepository{
			err: errors.New("network error"),
		}
		pricesRepo := &stubPricesRepository{
			err: errors.New("network error"),
		}

		// when
		etf := processETF("INVALID", dividendsRepo, pricesRepo)

		// then
		assert.Empty(t, etf.AmountDividendsPerYear)
		assert.Empty(t, etf.AverageClosingPricePerYear)
	})
}

func TestMain_ApplyColors(t *testing.T) {
	t.Parallel()

	t.Run("should apply green color to cells at or above the target yield and red below", func(t *testing.T) {
		t.Parallel()

		// given
		row := []string{"10.00%", "5.00%", "15.00%"}

		// when
		result := applyColors(row)

		// then
		expected := []string{
			ansiGreen + "10.00%" + ansiReset,
			ansiRed + "5.00%" + ansiReset,
			ansiGreen + "15.00%" + ansiReset,
		}
		assert.Equal(t, expected, result, "they should be equal")
	})

	t.Run("should leave non-percentage cells unchanged", func(t *testing.T) {
		t.Parallel()

		// given
		row := []string{"SPY", "$5.50", "N/A"}

		// when
		result := applyColors(row)

		// then
		assert.Equal(t, row, result)
	})
}
