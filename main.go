package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/olekukonko/tablewriter"
	"github.com/sirupsen/logrus"
)

// ETF represents an ETF and its dividend cash amounts by year
type ETF struct {
	Name          string
	DividendYears map[string]float64 // Key: Year, Value: Total Dividend Cash
}

// CalculateAverage calculates the average of the available dividend cash amounts
func (e *ETF) CalculateAverage() float64 {
	if len(e.DividendYears) == 0 {
		return 0
	}
	var sum float64
	for _, value := range e.DividendYears {
		sum += value
	}
	return sum / float64(len(e.DividendYears))
}

// DisplayYearlySums formats the yearly sums for table display
func (e *ETF) DisplayYearlySums(totalYears int, startYear int) []string {
	formatted := make([]string, totalYears)

	for i := 0; i < totalYears; i++ {
		year := fmt.Sprintf("%d", startYear-i)
		if value, exists := e.DividendYears[year]; exists {
			formatted[i] = fmt.Sprintf("$%.2f", value)
		} else {
			formatted[i] = "-"
		}
	}
	return formatted
}

// ScrapeETFs populates an array of ETFs with dividend cash data
func ScrapeETFs(etfNames []string) ([]*ETF, error) {
	var etfs []*ETF
	for _, name := range etfNames {
		etf := &ETF{Name: name, DividendYears: make(map[string]float64)}
		dividends, err := scrapeDividendCash(name)
		if err != nil {
			logrus.WithError(err).Errorf("Failed to scrape data for ETF: %s", name)
			continue
		}
		etf.DividendYears = dividends
		etfs = append(etfs, etf)
	}
	return etfs, nil
}

// scrapeDividendCash scrapes the total annual dividend cash amounts for the given ETF
func scrapeDividendCash(etf string) (map[string]float64, error) {
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
	err := c.Visit(url)
	if err != nil {
		return nil, err
	}

	return yearlyTotals, nil
}

func main() {
	logrus.Info("Starting ETF data scraping...")

	etfNames := []string{"HYGW", "RIET", "SDIV", "SVOL", "XYLD"}
	etfs, err := ScrapeETFs(etfNames)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to scrape ETFs")
	}

	table := tablewriter.NewWriter(os.Stdout)
	totalYears := 5
	currentYear := time.Now().Year()
	headers := []string{"ETF"}
	for i := 0; i < totalYears; i++ {
		headers = append(headers, fmt.Sprintf("%d", currentYear-i))
	}
	headers = append(headers, "Average")
	table.SetHeader(headers)

	for _, etf := range etfs {
		row := []string{etf.Name}
		row = append(row, etf.DisplayYearlySums(totalYears, currentYear)...)
		row = append(row, fmt.Sprintf("$%.2f", etf.CalculateAverage()))
		table.Append(row)
	}

	logrus.Info("Rendering the results...")
	table.Render()
}
