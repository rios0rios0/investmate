package main

import (
	"testing"

	"github.com/olekukonko/tablewriter"
	"github.com/stretchr/testify/assert"
)

func TestMain_ProcessETF(t *testing.T) {
	t.Parallel()

	t.Run("should process ETF data when a valid ETF name is given", func(t *testing.T) {
		t.Parallel()

		// given

		// when
		etf := processETF("SPY")

		// then
		assert.Equal(t, "SPY", etf.Name)
		assert.NotEmpty(t, etf.AmountDividendsPerYear)
		assert.NotEmpty(t, etf.AverageClosingPricePerYear)
	})

	t.Run("should return an empty ETF when an invalid ETF name is given", func(t *testing.T) {
		t.Parallel()

		// given

		// when
		etf := processETF("INVALID")

		// then
		assert.Empty(t, etf.AmountDividendsPerYear)
		assert.Empty(t, etf.AverageClosingPricePerYear)
	})
}

func TestMain_GetColors(t *testing.T) {
	t.Parallel()

	t.Run("should return colors for dividend yield row according the target yield", func(t *testing.T) {
		t.Parallel()

		// given
		row := []string{"10.00%", "5.00%", "15.00%"}

		// when
		result := getColors(row)

		// then
		expected := []tablewriter.Colors{
			{tablewriter.FgGreenColor},
			{tablewriter.FgRedColor},
			{tablewriter.FgGreenColor},
		}
		assert.Equal(t, expected, result, "they should be equal")
	})
}
