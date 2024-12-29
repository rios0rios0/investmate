package entities_test

import (
	"testing"

	"github.com/rios0rios0/investmate/internal/domain/entities"
	"github.com/stretchr/testify/suite"
)

type ETFTestSuite struct {
	suite.Suite
	etf *entities.ETF
}

func (suite *ETFTestSuite) SetupTest() {
	suite.etf = &entities.ETF{
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
		// given
		// on the setup

		// when
		result := suite.etf.ShowDividendsPerYear(2023, 5)

		// then
		expected := []string{"$10.000", "$15.000", "$20.000", "-", "-"}
		suite.Equal(expected, result)
	})
}

func (suite *ETFTestSuite) TestAverageDividends() {
	suite.Run("should calculate average dividends", func() {
		// given
		// on the setup

		// when
		result := suite.etf.AverageDividends(2023, 5)

		// then
		expected := 15.0
		suite.InEpsilon(expected, result, 0.01)
	})
}

func (suite *ETFTestSuite) TestShowClosingPricesPerYear() {
	suite.Run("should return formatted closing prices per year", func() {
		// given
		// on the setup

		// when
		result := suite.etf.ShowClosingPricesPerYear(2023, 5)

		// then
		expected := []string{"$100.000", "$150.000", "$200.000", "-", "-"}
		suite.Equal(expected, result)
	})
}

func (suite *ETFTestSuite) TestAverageClosingPrices() {
	suite.Run("should calculate average closing prices", func() {
		// given
		// on the setup

		// when
		result := suite.etf.AverageClosingPrices(2023, 5)

		// then
		expected := 150.0
		suite.InEpsilon(expected, result, 0.01)
	})
}

func (suite *ETFTestSuite) TestShowDividendYieldPerYear() {
	suite.Run("should return formatted dividend yields per year", func() {
		// given
		// on the setup

		// when
		result := suite.etf.ShowDividendYieldPerYear(2023, 5)

		// then
		expected := []string{"10.000%", "10.000%", "10.000%", "-", "-"}
		suite.Equal(expected, result)
	})
}

func (suite *ETFTestSuite) TestAverageDividendYield() {
	suite.Run("should calculate average dividend yield", func() {
		// given

		// when
		result := suite.etf.AverageDividendYield(2023, 5)

		// then
		expected := 15.0
		suite.InEpsilon(expected, result, 0.01)
	})
}

func TestETFTestSuite(t *testing.T) {
	suite.Run(t, new(ETFTestSuite))
}
