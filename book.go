package gdax

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

// BidAsk ...
type BidAsk struct {
	Price          float64
	Size           float64
	NumberOfOrders int
	OrderID        string
}

// Bids ...
type Bids []BidAsk

// Asks ...
type Asks []BidAsk

// Book represents the Bid and Ask
type Book struct {
	Sequence int `json:"sequence"`
	Bids     `json:"bids"`
	Asks     `json:"asks"`
}

// UnmarshalJSON takes the array and assigns the values to the BidAsk Struct
func (b *BidAsk) UnmarshalJSON(data []byte) error {
	var entry []interface{}

	if err := json.Unmarshal(data, &entry); err != nil {
		return err
	}

	price, err := strconv.ParseFloat(entry[0].(string), 32)
	if err != nil {
		return err
	}

	size, err := strconv.ParseFloat(entry[1].(string), 32)
	if err != nil {
		return err
	}

	*b = BidAsk{
		Price: price,
		Size:  size,
	}

	var orderID string
	numOfOrders, ok := entry[2].(float64)
	if !ok {
		orderID, ok = entry[2].(string)
		if !ok {
			return errors.New("Could not unmarshal last column on BidAsk")
		}
		b.OrderID = orderID
	} else {
		b.NumberOfOrders = int(numOfOrders)
	}
	return nil
}

// BookClient API client
type BookClient struct {
	client *Client
}

// GetBookByProduct ...
func (bc BookClient) GetBookByProduct(productID string) (Book, error) {
	var book Book
	_, err := bc.client.Request("GET", fmt.Sprintf("/products/%s/book", productID), nil, &book)
	return book, err
}
