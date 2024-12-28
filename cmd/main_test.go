package main

import (
	"testing"

	"github.com/olekukonko/tablewriter"
	"github.com/stretchr/testify/suite"
)

type ETFTestSuite struct {
	suite.Suite
	etf *ETF
}

func (suite *ETFTestSuite) SetupTest() {
	suite.etf = &ETF{
		Name: "TestETF",
		AmountDividendsPerYear: map[string]float64{
			"2023": 10.0,
			"2022": 15.0,
			"2021": 20.0,
		},
		AverageClosingPricePerYear: map[string]float64{
			"2023": 100.0,
			"2022": 150.0,
			"2021": 200.0,
		},
		DividendYieldPerYear: map[string]float64{
			"2023": 10.0,
			"2022": 15.0,
			"2021": 20.0,
		},
	}
}

func (suite *ETFTestSuite) TestShowDividendsPerYear() {
	suite.Run("should return formatted dividends per year", func() {
		expected := []string{"$10.00", "$15.00", "$20.00", "-", "-"}
		result := suite.etf.ShowDividendsPerYear(2023, 5)
		suite.Equal(expected, result)
	})
}

func (suite *ETFTestSuite) TestAverageDividends() {
	suite.Run("should calculate average dividends", func() {
		expected := 15.0
		result := suite.etf.AverageDividends(2023, 5)
		suite.InEpsilon(expected, result, 0.01)
	})
}

func (suite *ETFTestSuite) TestShowClosingPricesPerYear() {
	suite.Run("should return formatted closing prices per year", func() {
		expected := []string{"$100.00", "$150.00", "$200.00", "-", "-"}
		result := suite.etf.ShowClosingPricesPerYear(2023, 5)
		suite.Equal(expected, result)
	})
}

func (suite *ETFTestSuite) TestAverageClosingPrices() {
	suite.Run("should calculate average closing prices", func() {
		expected := 150.0
		result := suite.etf.AverageClosingPrices(2023, 5)
		suite.InEpsilon(expected, result, 0.01)
	})
}

func (suite *ETFTestSuite) TestShowDividendYieldPerYear() {
	suite.Run("should return formatted dividend yields per year", func() {
		expected := []string{"10.00%", "10.00%", "10.00%", "-", "-"}
		result := suite.etf.ShowDividendYieldPerYear(2023, 5)
		suite.Equal(expected, result)
	})
}

func (suite *ETFTestSuite) TestAverageDividendYield() {
	suite.Run("should calculate average dividend yield", func() {
		expected := 15.0
		result := suite.etf.AverageDividendYield(2023, 5)
		suite.InEpsilon(expected, result, 0.01)
	})
}

func (suite *ETFTestSuite) TestProcessETF() {
	suite.Run("should process ETF data", func() {
		etf := processETF("SPY")
		suite.Equal("SPY", etf.Name)
		suite.NotEmpty(etf.AmountDividendsPerYear)
		suite.NotEmpty(etf.AverageClosingPricePerYear)
	})
}

func (suite *ETFTestSuite) TestGetColors() {
	suite.Run("should return colors for dividend yield row", func() {
		row := []string{"10.00%", "5.00%", "15.00%"}
		expected := []tablewriter.Colors{
			{tablewriter.FgGreenColor},
			{tablewriter.FgRedColor},
			{tablewriter.FgGreenColor},
		}
		result := getColors(row)
		suite.Equal(expected, result)
	})
}

func TestETFTestSuite(t *testing.T) {
	suite.Run(t, new(ETFTestSuite))
}
