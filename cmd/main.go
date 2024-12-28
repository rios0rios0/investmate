package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/olekukonko/tablewriter"
	"github.com/sirupsen/logrus"
)

const (
	YearsToFetch         = 5 // Number of years to fetch data for
	PercentageMultiplier = 100
	NumberOfDaysInYear   = 365
)

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
			formatted[i] = fmt.Sprintf("$%.2f", value)
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
			formatted[i] = fmt.Sprintf("$%.2f", value)
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
				formatted[i] = fmt.Sprintf("%.2f%%", yield)
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

// processETF populates an array of ETFs with dividend cash data and average closing prices
func processETF(name string) *ETF {
	etf := &ETF{
		Name:                       name,
		AmountDividendsPerYear:     make(map[string]float64),
		AverageClosingPricePerYear: make(map[string]float64),
	}

	dividendsPerYear, err := crawlingDividendsPerYear(name)
	if err != nil {
		logrus.WithError(err).Errorf("Failed to scrape data for ETF: %s", name)
	}
	etf.AmountDividendsPerYear = dividendsPerYear

	closingPricesPerYear, err := fetchAverageClosingPricesPerYear(name)
	if err != nil {
		logrus.WithError(err).Errorf("Failed to fetch average close prices for ETF: %s", name)
	}
	etf.AverageClosingPricePerYear = closingPricesPerYear

	return etf
}

// crawlingDividendsPerYear scrapes the total annual dividend cash amounts for the given ETF
func crawlingDividendsPerYear(etf string) (map[string]float64, error) {
	c := colly.NewCollector()

	yearlyTotals := make(map[string]float64)

	c.OnHTML("table#dividend_table tbody tr", func(e *colly.HTMLElement) {
		date := e.ChildText("td:nth-child(2)")        // Payout Date
		dividendStr := e.ChildText("td:nth-child(3)") // Cash Amount

		year := strings.Split(date, "-")[0]

		dividendStr = strings.ReplaceAll(dividendStr, "$", "")
		dividendStr = strings.TrimSpace(dividendStr)
		dividend, err := strconv.ParseFloat(dividendStr, 64)
		if err == nil && year != "" {
			yearlyTotals[year] += dividend
		}
	})

	url := fmt.Sprintf("https://dividendhistory.org/payout/%s/", etf)
	if err := c.Visit(url); err != nil {
		return nil, fmt.Errorf("failed to visit URL: %w", err)
	}

	return yearlyTotals, nil
}

// fetchAverageClosingPricesPerYear fetches the average closing prices per year for the given ETF
func fetchAverageClosingPricesPerYear(etf string) (map[string]float64, error) {
	averageClosePrices := make(map[string]float64)
	currentYear := time.Now().Year()
	fromDate := fmt.Sprintf("%d-01-01", currentYear-YearsToFetch)
	toDate := fmt.Sprintf("%d-12-31", currentYear)

	url := fmt.Sprintf(
		"https://api.nasdaq.com/api/quote/%s/historical?assetclass=etf&fromdate=%s&todate=%s&limit=%d&offset=0",
		etf, fromDate, toDate, YearsToFetch*NumberOfDaysInYear,
	)
	ctx := context.Background()
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	req.Header.Set(
		"User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 "+
			"(KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36 Edg/131.0.0.0",
	)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		Data struct {
			TradesTable struct {
				Rows []struct {
					Close string `json:"close"`
					Date  string `json:"date"`
				} `json:"rows"`
			} `json:"tradesTable"`
		} `json:"data"`
	}

	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	yearlySums := make(map[string]float64)
	yearlyCounts := make(map[string]int)

	for _, row := range result.Data.TradesTable.Rows {
		closePrice, parseErr := strconv.ParseFloat(strings.ReplaceAll(row.Close, "$", ""), 64)
		if parseErr == nil {
			year := strings.Split(row.Date, "/")[2]
			yearlySums[year] += closePrice
			yearlyCounts[year]++
		}
	}

	for year, sum := range yearlySums {
		if count, exists := yearlyCounts[year]; exists && count > 0 {
			averageClosePrices[year] = sum / float64(count)
		}
	}

	return averageClosePrices, nil
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
	logrus.Info("Starting ETF data scraping...")

	var etfs []*ETF
	etfNames := []string{"HYGW", "RIET", "SDIV", "SVOL", "XYLD"}

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
		dividendRow = append(dividendRow, fmt.Sprintf("$%.2f", etf.AverageDividends(currentYear, totalYears)))
		table.Append(dividendRow)

		// Closing prices
		closePriceRow := []string{etf.Name + " Closing Prices"}
		closePriceRow = append(closePriceRow, etf.ShowClosingPricesPerYear(currentYear, totalYears)...)
		closePriceRow = append(closePriceRow, fmt.Sprintf("$%.2f", etf.AverageClosingPrices(currentYear, totalYears)))
		table.Append(closePriceRow)

		// Dividend yields
		dividendYieldRow := []string{etf.Name + " Dividend Yields"}
		dividendYieldRow = append(dividendYieldRow, etf.ShowDividendYieldPerYear(currentYear, totalYears)...)
		dividendYieldRow = append(dividendYieldRow, fmt.Sprintf("%.2f%%", etf.AverageDividendYield(currentYear, totalYears)))
		table.Rich(dividendYieldRow, getColors(dividendYieldRow))

		// Add a separator row after every 3 lines
		table.Append([]string{"-", "-", "-", "-", "-", "-", "-"})
	}

	logrus.Info("Rendering the results...")
	table.Render()
}
