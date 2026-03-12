package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/rios0rios0/investmate/internal/domain/entities"
	"github.com/rios0rios0/investmate/internal/domain/repositories"
	"github.com/rios0rios0/investmate/internal/infrastructure/repositories/nasdaq"
	logger "github.com/sirupsen/logrus"
)

const (
	// YearsToFetch is the number of years to fetch data for.
	YearsToFetch = 5

	// targetYieldPercentage is the minimum dividend yield percentage considered a good target.
	targetYieldPercentage = 9

	// ansiGreen is the ANSI escape code for green foreground text.
	ansiGreen = "\033[32m"

	// ansiRed is the ANSI escape code for red foreground text.
	ansiRed = "\033[31m"

	// ansiReset is the ANSI escape code to reset text formatting.
	ansiReset = "\033[0m"
)

// processETF populates an ETF struct with dividend cash data and average closing prices.
func processETF(
	name string,
	dividendsRepo repositories.DividendsRepository,
	pricesRepo repositories.PricesRepository,
) *entities.ETF {
	etf := &entities.ETF{
		Name:                       name,
		AmountDividendsPerYear:     make(map[string]float64),
		AverageClosingPricePerYear: make(map[string]float64),
	}

	dividendsPerYear, err := dividendsRepo.ListDividendsByETF(name)
	if err != nil {
		logger.WithError(err).Errorf("Failed to scrape data for ETF: %s", name)
	}

	etf.AmountDividendsPerYear = dividendsPerYear

	closingPricesPerYear, err := pricesRepo.ListClosingPricesByETF(name)
	if err != nil {
		logger.WithError(err).Errorf("Failed to fetch average close prices for ETF: %s", name)
	}

	etf.AverageClosingPricePerYear = closingPricesPerYear

	return etf
}

// applyColors wraps each cell that contains a percentage value with the appropriate ANSI color code.
// Cells with a dividend yield at or above the target threshold are colored green; below is red.
func applyColors(row []string) []string {
	colored := make([]string, len(row))

	for i, cell := range row {
		before, ok := strings.CutSuffix(cell, "%")
		if !ok {
			colored[i] = cell
			continue
		}

		value, err := strconv.ParseFloat(before, 64)

		switch {
		case err == nil && value >= targetYieldPercentage:
			colored[i] = ansiGreen + cell + ansiReset
		case err == nil && value < targetYieldPercentage:
			colored[i] = ansiRed + cell + ansiReset
		default:
			colored[i] = cell
		}
	}

	return colored
}

func main() {
	logger.Info("Starting ETF data scraping...")

	var etfs []*entities.ETF
	etfNames := []string{
		"SPY", "QQQ", "SCHD", "YYY", "GLD",
		"HYGW", "RIET", "SDIV", "SVOL", "XYLD",
	}

	dividendsRepo := nasdaq.NewAPIDividendsRepository()
	pricesRepo := nasdaq.NewAPIPricesRepository()

	for _, name := range etfNames {
		etf := processETF(name, dividendsRepo, pricesRepo)
		etfs = append(etfs, etf)
	}

	table := tablewriter.NewWriter(os.Stdout)
	totalYears := YearsToFetch
	currentYear := time.Now().Year()
	headers := []string{"ETF"}

	for i := range totalYears {
		headers = append(headers, strconv.Itoa(currentYear-i))
	}

	headers = append(headers,
		"Averages",
		"Payout Frequency", "Average Volume", "Expense Ratio", "Beta", "AUM", "Inception Date",
	)
	table.Header(headers)

	for _, etf := range etfs {
		// Dividend sums.
		dividendRow := []string{etf.Name + " Dividends"}
		dividendRow = append(dividendRow, etf.ShowDividendsPerYear(currentYear, totalYears)...)
		dividendRow = append(dividendRow, fmt.Sprintf("$%.3f", etf.AverageDividends(currentYear, totalYears)))

		if err := table.Append(dividendRow); err != nil {
			logger.WithError(err).Errorf("Failed to append dividend row for ETF: %s", etf.Name)
		}

		// Closing prices.
		closePriceRow := []string{etf.Name + " Closing Prices"}
		closePriceRow = append(closePriceRow, etf.ShowClosingPricesPerYear(currentYear, totalYears)...)
		closePriceRow = append(closePriceRow, fmt.Sprintf("$%.3f", etf.AverageClosingPrices(currentYear, totalYears)))

		if err := table.Append(closePriceRow); err != nil {
			logger.WithError(err).Errorf("Failed to append close price row for ETF: %s", etf.Name)
		}

		// Dividend yields with color-coded cells based on the target yield threshold.
		dividendYieldRow := []string{etf.Name + " Dividend Yields"}
		dividendYieldRow = append(dividendYieldRow, etf.ShowDividendYieldPerYear(currentYear, totalYears)...)
		dividendYieldRow = append(
			dividendYieldRow,
			fmt.Sprintf("%.3f%%", etf.AverageDividendYield(currentYear, totalYears)),
		)

		if err := table.Append(applyColors(dividendYieldRow)); err != nil {
			logger.WithError(err).Errorf("Failed to append dividend yield row for ETF: %s", etf.Name)
		}

		// Add a separator row after every 3 lines.
		if err := table.Append([]string{"-", "-", "-", "-", "-", "-", "-"}); err != nil {
			logger.WithError(err).Error("Failed to append separator row")
		}
	}

	logger.Info("Rendering the results...")

	if err := table.Render(); err != nil {
		logger.WithError(err).Fatal("Failed to render the table")
	}
}
