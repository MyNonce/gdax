package gdax

import (
	"errors"
	"testing"
)

func TestCreateLimitOrders(t *testing.T) {
	order := Order{
		Price:     9315,
		Size:      .01,
		Side:      "buy",
		ProductID: "BTC-USD",
	}

	savedOrder, err := GDAXClient.OrderClient.CreateOrder(order)
	if err != nil {
		t.Error(err)
	}

	if savedOrder.ID == "" {
		t.Error(errors.New("No create id found"))
	}

	if err := GDAXClient.OrderClient.CancelOrder(savedOrder.ID); err != nil {
		t.Error(err)
	}
}
