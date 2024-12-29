package nasdaq

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	YearsToFetch       = 5 // Number of years to fetch data for
	NumberOfDaysInYear = 365
)

type APIPricesRepository struct {
}

func NewAPIPricesRepository() *APIPricesRepository {
	return &APIPricesRepository{}
}

func (r APIPricesRepository) ListClosingPricesByETF(etf string) (map[string]float64, error) {
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
