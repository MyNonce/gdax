package gdax

import "time"

// Fill represents orders that have been filled
type Fill struct {
	TradeID   int       `json:"trade_id,int"`
	ProductID string    `json:"product_id"`
	Price     float64   `json:"price,string"`
	Size      float64   `json:"size,string"`
	FillID    string    `json:"order_id"`
	CreatedAt time.Time `json:"created_at,string"`
	Fee       float64   `json:"fee,string"`
	Settled   bool      `json:"settled"`
	Side      string    `json:"side"`
	Liquidity string    `json:"liquidity"`
}

// SumOrder returns the total recieved amount
func (f Fill) SumOrder() float64 {
	return f.Size*f.Price - f.Fee
}

// FillsClient API client
type FillsClient struct {
	client *Client
}

func addExtraParam(key, value string, p *PaginationParams) {
	if (p == nil) || (key == "") || (value == "") {
		return
	}
	p.AddExtraParam(key, value)
}

// ListFills pulls
func (fc *FillsClient) ListFills(orderID, productID string) *Cursor {
	paginationParams := PaginationParams{}
	addExtraParam("order_id", orderID, &paginationParams)
	addExtraParam("product_id", productID, &paginationParams)
	return NewCursor(fc.client, "GET", "/fills", &paginationParams)
}
