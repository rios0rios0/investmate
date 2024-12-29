package entities

import (
	"fmt"
	"strconv"
)

const PercentageMultiplier = 100

// ETF represents an ETF and its dividend cash amounts by year
type ETF struct {
	Name                       string
	AmountDividendsPerYear     map[string]float64 // Key: Year, Value: Total Dividend Cash
	AverageClosingPricePerYear map[string]float64 // Key: Year, Value: Average Closing Price
	DividendYieldPerYear       map[string]float64 // Key: Year, Value: Dividend Yield Percentage
}

// ShowDividendsPerYear formats the yearly sums for table display
func (e *ETF) ShowDividendsPerYear(startYear, totalYears int) []string {
	formatted := make([]string, totalYears)

	for i := range make([]struct{}, totalYears) {
		year := strconv.Itoa(startYear - i)
		if value, exists := e.AmountDividendsPerYear[year]; exists {
			formatted[i] = fmt.Sprintf("$%.3f", value)
		} else {
			formatted[i] = "-"
		}
	}
	return formatted
}

// AverageDividends calculates the average of the available dividend cash amounts for the specified years
func (e *ETF) AverageDividends(startYear, totalYears int) float64 {
	if len(e.AmountDividendsPerYear) == 0 {
		return 0
	}

	var sum float64
	var count int
	for i := range make([]struct{}, totalYears) {
		year := strconv.Itoa(startYear - i)
		if value, exists := e.AmountDividendsPerYear[year]; exists {
			sum += value
			count++
		}
	}
	if count == 0 {
		return 0
	}

	return sum / float64(count)
}

// ShowClosingPricesPerYear formats the average closing prices for table display
func (e *ETF) ShowClosingPricesPerYear(startYear, totalYears int) []string {
	formatted := make([]string, totalYears)

	for i := range make([]struct{}, totalYears) {
		year := strconv.Itoa(startYear - i)
		if value, exists := e.AverageClosingPricePerYear[year]; exists {
			formatted[i] = fmt.Sprintf("$%.3f", value)
		} else {
			formatted[i] = "-"
		}
	}
	return formatted
}

// AverageClosingPrices calculates the average closing prices for the specified years
func (e *ETF) AverageClosingPrices(startYear, totalYears int) float64 {
	if len(e.AverageClosingPricePerYear) == 0 {
		return 0
	}

	var sum float64
	var count int
	for i := range make([]struct{}, totalYears) {
		year := strconv.Itoa(startYear - i)
		if value, exists := e.AverageClosingPricePerYear[year]; exists {
			sum += value
			count++
		}
	}
	if count == 0 {
		return 0
	}

	return sum / float64(count)
}

// ShowDividendYieldPerYear calculates the dividend yield for each year and stores it in the ETF struct
func (e *ETF) ShowDividendYieldPerYear(startYear, totalYears int) []string {
	formatted := make([]string, totalYears)
	e.DividendYieldPerYear = make(map[string]float64)

	for i := range make([]struct{}, totalYears) {
		year := strconv.Itoa(startYear - i)
		if dividend, dividendExists := e.AmountDividendsPerYear[year]; dividendExists {
			if closingPrice, priceExists := e.AverageClosingPricePerYear[year]; priceExists && closingPrice != 0 {
				yield := (dividend / closingPrice) * PercentageMultiplier
				e.DividendYieldPerYear[year] = yield
				formatted[i] = fmt.Sprintf("%.3f%%", yield)
			} else {
				formatted[i] = "-"
			}
		} else {
			formatted[i] = "-"
		}
	}

	return formatted
}

// AverageDividendYield calculates the average dividend yield for the specified years
func (e *ETF) AverageDividendYield(startYear int, totalYears int) float64 {
	if len(e.DividendYieldPerYear) == 0 {
		return 0
	}

	var sum float64
	var count int
	for i := range make([]struct{}, totalYears) {
		year := strconv.Itoa(startYear - i)
		if yield, exists := e.DividendYieldPerYear[year]; exists {
			sum += yield
			count++
		}
	}
	if count == 0 {
		return 0
	}

	return sum / float64(count)
}
