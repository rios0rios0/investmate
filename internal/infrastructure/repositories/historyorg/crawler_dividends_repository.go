package historyorg

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

type CrawlerDividendsRepository struct {
}

func NewCrawlerDividendsRepository() *CrawlerDividendsRepository {
	return &CrawlerDividendsRepository{}
}

func (r CrawlerDividendsRepository) ListDividendsByETF(etf string) (map[string]float64, error) {
	c := colly.NewCollector()

	yearlyTotals := make(map[string]float64)

	c.OnHTML("table#dividend_table tbody tr", func(e *colly.HTMLElement) {
		date := e.ChildText("td:nth-child(2)")        // Payout Date
		dividendStr := e.ChildText("td:nth-child(3)") // Cash Amount

		year := strings.Split(date, "-")[0]

		dividendStr = strings.TrimSpace(strings.ReplaceAll(dividendStr, "$", ""))
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
