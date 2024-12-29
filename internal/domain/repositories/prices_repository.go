package repositories

// PricesRepository defines the interface for getting prices per year.
type PricesRepository interface {
	ListClosingPricesByETF(etf string) (map[string]float64, error)
}
