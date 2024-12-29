package nasdaq

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type APIDividendsRepository struct {
}

func NewAPIDividendsRepository() *APIDividendsRepository {
	return &APIDividendsRepository{}
}

func (r *APIDividendsRepository) ListDividendsByETF(etf string) (map[string]float64, error) {
	url := fmt.Sprintf("https://api.nasdaq.com/api/quote/%s/dividends?assetclass=etf", etf)

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
			Dividends struct {
				Rows []struct {
					Amount      string `json:"amount"`
					PaymentDate string `json:"paymentDate"`
				} `json:"rows"`
			} `json:"dividends"`
		} `json:"data"`
	}

	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	yearlySums := make(map[string]float64)

	for _, row := range result.Data.Dividends.Rows {
		amount, parseErr := strconv.ParseFloat(strings.ReplaceAll(row.Amount, "$", ""), 64)
		if parseErr == nil {
			year := strings.Split(row.PaymentDate, "/")[2]
			yearlySums[year] += amount
		}
	}

	return yearlySums, nil
}
