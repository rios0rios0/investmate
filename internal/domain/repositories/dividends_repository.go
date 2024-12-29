package repositories

// DividendsRepository defines the interface for getting dividends per year.
type DividendsRepository interface {
	ListDividendsByETF(etf string) (map[string]float64, error)
}
