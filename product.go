package gdax

// Product represents the available currencty pairs for trading on GDAX
type Product struct {
	ID             string  `json:"id"`
	BaseCurrency   string  `json:"base_currency"`
	QuoteCurrency  string  `json:"quote_currency"`
	BaseMinSize    float64 `json:"base_min_size"`
	BaseMaxSize    float64 `json:"base_max_size"`
	QuoteIncrement float64 `json:"quote_increment"`
}

// ProductClient API client
type ProductClient struct {
	client *Client
}

// GetProducts ...
func (pc ProductClient) GetProducts() ([]Product, error) {
	var products []Product
	_, err := pc.client.Request("GET", "/products", nil, &products)
	return products, err
}
