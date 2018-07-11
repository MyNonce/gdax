package gdax

import (
	"fmt"
	"time"
)
	
// Order ...
type Order struct {
	Type      string  `json:"type"`
	Size      float64 `json:"size,string,omitempty"`
	Side      string  `json:"side"`
	ProductID string  `json:"product_id"`
	ClientOID string  `json:"client_oid,omitempty"`
	Stp       string  `json:"stp,omitempty"`
	// Limit Order
	Price       float64 `json:"price,string,omitempty"`
	TimeInForce string  `json:"time_in_force,omitempty"`
	PostOnly    bool    `json:"post_only,omitempty"`
	CancelAfter string  `json:"cancel_after,omitempty"`
	// Market Order
	Funds float64 `json:"funds,string,omitempty"`
	// Response Fields
	ID            string     `json:"id"`
	Status        string     `json:"status,omitempty"`
	Settled       bool       `json:"settled,omitempty"`
	DoneReason    string     `json:"done_reason,omitempty"`
	CreatedAt     *time.Time `json:"created_at,string,omitempty"`
	FillFees      float64    `json:"fill_fees,string,omitempty"`
	FilledSize    float64    `json:"filled_size,string,omitempty"`
	ExecutedValue float64    `json:"executed_value,string,omitempty"`
}

// NewLimitSellOrder ...
func NewLimitSellOrder(productID string, p, s float64) *Order {
	o := &Order{
		Side:      "sell",
		Type:      "limit",
		ProductID: productID,
		Price:     p,
		Size:      s,
		PostOnly:  true,
	}
	return o
}

// NewLimitBuyOrder ...
func NewLimitBuyOrder(productID string, p, s float64) *Order {
	o := &Order{
		Side:        "buy",
		Type:        "limit",
		ProductID:   productID,
		Price:       p,
		Size:        s,
		PostOnly:    true,
		TimeInForce: "GTT",
		CancelAfter: "min",
	}
	return o
}

// OrderClient ...
type OrderClient struct {
	client *Client
}

// GetOrdersByID ...
func (oc OrderClient) GetOrdersByID(orderID string) (Order, error) {
	var order Order
	_, err := oc.client.Request("GET", fmt.Sprintf("/orders/%s", orderID), nil, &order)
	return order, err
}

// CreateOrder ...
func (oc OrderClient) CreateOrder(newOrder Order) (Order, error) {
	var order Order
	_, err := oc.client.Request("POST", "/orders", &newOrder, &order)
	return order, err
}

// CancelOrder ...
func (oc OrderClient) CancelOrder(orderID string) error {
	_, err := oc.client.Request("DELETE", fmt.Sprintf("/orders/%s", orderID), nil, nil)
	return err
}

// CancelAllOrders ...
func (oc OrderClient) CancelAllOrders() ([]string, error) {
	var orders []string
	_, err := oc.client.Request("DELETE", "/orders", nil, orders)
	return orders, err
}
