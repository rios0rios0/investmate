package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/rios0rios0/investmate/internal/domain/entities"
	"github.com/rios0rios0/investmate/internal/infrastructure/repositories/nasdaq"
	logger "github.com/sirupsen/logrus"
)

const (
	YearsToFetch = 5 // Number of years to fetch data for
)

// processETF populates an array of ETFs with dividend cash data and average closing prices
func processETF(name string) *entities.ETF {
	etf := &entities.ETF{
		Name:                       name,
		AmountDividendsPerYear:     make(map[string]float64),
		AverageClosingPricePerYear: make(map[string]float64),
	}

	// dividendsPerYear, err := statusinvest.NewCrawlerDividendsRepository().ListDividendsByETF(name)
	dividendsPerYear, err := nasdaq.NewAPIDividendsRepository().ListDividendsByETF(name)
	if err != nil {
		logger.WithError(err).Errorf("Failed to scrape data for ETF: %s", name)
	}
	etf.AmountDividendsPerYear = dividendsPerYear

	closingPricesPerYear, err := nasdaq.NewAPIPricesRepository().ListClosingPricesByETF(name)
	if err != nil {
		logger.WithError(err).Errorf("Failed to fetch average close prices for ETF: %s", name)
	}
	etf.AverageClosingPricePerYear = closingPricesPerYear

	return etf
}

// getColors returns the colors for the dividend yield row
func getColors(row []string) []tablewriter.Colors {
	colors := make([]tablewriter.Colors, len(row))
	for i, cell := range row {
		if strings.HasSuffix(cell, "%") {
			value, err := strconv.ParseFloat(strings.TrimSuffix(cell, "%"), 64)
			if err == nil && value >= 9 {
				colors[i] = tablewriter.Colors{tablewriter.FgGreenColor}
			} else if err == nil && value < 9 {
				colors[i] = tablewriter.Colors{tablewriter.FgRedColor}
			}
		}
	}
	return colors
}

func main() {
	logger.Info("Starting ETF data scraping...")

	var etfs []*entities.ETF
	etfNames := []string{
		"SPY", "QQQ", "SCHD", "YYY", "GLD",
		"HYGW", "RIET", "SDIV", "SVOL", "XYLD",
	}

	for _, name := range etfNames {
		etf := processETF(name)
		etfs = append(etfs, etf)
	}

	table := tablewriter.NewWriter(os.Stdout)
	totalYears := YearsToFetch
	currentYear := time.Now().Year()
	headers := []string{"ETF"}
	for i := range make([]struct{}, totalYears) {
		headers = append(headers, strconv.Itoa(currentYear-i))
	}
	headers = append(headers, "Averages")
	table.SetHeader(headers)

	for _, etf := range etfs {
		// Dividend sums
		dividendRow := []string{etf.Name + " Dividends"}
		dividendRow = append(dividendRow, etf.ShowDividendsPerYear(currentYear, totalYears)...)
		dividendRow = append(dividendRow, fmt.Sprintf("$%.3f", etf.AverageDividends(currentYear, totalYears)))
		table.Append(dividendRow)

		// Closing prices
		closePriceRow := []string{etf.Name + " Closing Prices"}
		closePriceRow = append(closePriceRow, etf.ShowClosingPricesPerYear(currentYear, totalYears)...)
		closePriceRow = append(closePriceRow, fmt.Sprintf("$%.3f", etf.AverageClosingPrices(currentYear, totalYears)))
		table.Append(closePriceRow)

		// Dividend yields
		dividendYieldRow := []string{etf.Name + " Dividend Yields"}
		dividendYieldRow = append(dividendYieldRow, etf.ShowDividendYieldPerYear(currentYear, totalYears)...)
		dividendYieldRow = append(dividendYieldRow, fmt.Sprintf("%.3f%%", etf.AverageDividendYield(currentYear, totalYears)))
		table.Rich(dividendYieldRow, getColors(dividendYieldRow))

		// Add a separator row after every 3 lines
		table.Append([]string{"-", "-", "-", "-", "-", "-", "-"})
	}

	logger.Info("Rendering the results...")
	table.Render()
}
