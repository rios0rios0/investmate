package nasdaq

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type APIDividendsRepository struct {
	client *APIClient
}

func NewAPIDividendsRepository() *APIDividendsRepository {
	return NewAPIDividendsRepositoryWithClient(NewAPIClient(DefaultBaseURL, http.DefaultClient))
}

// NewAPIDividendsRepositoryWithClient builds the repository on top of an explicit APIClient, which is how
// tests point it at an in-memory server. A nil client falls back to one bound to [DefaultBaseURL].
func NewAPIDividendsRepositoryWithClient(client *APIClient) *APIDividendsRepository {
	return &APIDividendsRepository{client: client}
}

func (r *APIDividendsRepository) ListDividendsByETF(etf string) (map[string]float64, error) {
	endpoint := fmt.Sprintf("/api/quote/%s/dividends?assetclass=etf", etf)

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

	if err := clientOrDefault(r.client).FetchJSON(context.Background(), endpoint, &result); err != nil {
		return nil, err
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
