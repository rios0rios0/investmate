package statusinvest

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gocolly/colly"
	logger "github.com/sirupsen/logrus"
)

type CrawlerDividendsRepository struct {
}

func NewCrawlerDividendsRepository() *CrawlerDividendsRepository {
	return &CrawlerDividendsRepository{}
}

func (r CrawlerDividendsRepository) ListDividendsByETF(etf string) (map[string]float64, error) {
	c := colly.NewCollector()

	yearlyTotals := make(map[string]float64)

	c.OnHTML("div#earning-section input#results", func(e *colly.HTMLElement) {
		jsonData := e.Attr("value")

		var dividends []struct {
			Value       float64 `json:"v"`
			PaymentDate string  `json:"pd"`
		}

		if err := json.Unmarshal([]byte(jsonData), &dividends); err != nil {
			logger.Errorf("Failed to unmarshal JSON data: %v", err)
			return
		}

		for _, dividend := range dividends {
			year := dividend.PaymentDate
			if year != "-" {
				year = strings.Split(dividend.PaymentDate, "/")[2]
			}
			yearlyTotals[year] += dividend.Value
		}
	})

	url := fmt.Sprintf("https://statusinvest.com.br/etf/eua/%s", etf)
	if err := c.Visit(url); err != nil {
		return nil, fmt.Errorf("failed to visit URL: %w", err)
	}

	return yearlyTotals, nil
}
