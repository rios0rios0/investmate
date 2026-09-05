package nasdaq

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	// YearsToFetch is the number of years to fetch data for.
	YearsToFetch = 5

	// NumberOfDaysInYear is the approximate number of trading data points in a year.
	NumberOfDaysInYear = 365
)

type APIPricesRepository struct {
	client *APIClient
}

func NewAPIPricesRepository() *APIPricesRepository {
	return NewAPIPricesRepositoryWithClient(NewAPIClient(DefaultBaseURL, http.DefaultClient))
}

// NewAPIPricesRepositoryWithClient builds the repository on top of an explicit APIClient, which is how
// tests point it at an in-memory server.
func NewAPIPricesRepositoryWithClient(client *APIClient) *APIPricesRepository {
	return &APIPricesRepository{client: client}
}

func (r APIPricesRepository) ListClosingPricesByETF(etf string) (map[string]float64, error) {
	averageClosePrices := make(map[string]float64)
	currentYear := time.Now().Year()
	fromDate := fmt.Sprintf("%d-01-01", currentYear-YearsToFetch)
	toDate := fmt.Sprintf("%d-12-31", currentYear)

	endpoint := fmt.Sprintf(
		"/api/quote/%s/historical?assetclass=etf&fromdate=%s&todate=%s&limit=%d&offset=0",
		etf, fromDate, toDate, YearsToFetch*NumberOfDaysInYear,
	)

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

	if err := r.client.FetchJSON(context.Background(), endpoint, &result); err != nil {
		return nil, err
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
