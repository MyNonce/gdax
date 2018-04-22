package gdax

import (
	"fmt"
	"log"
	"time"
)

// Ledger represents account changes
type Ledger struct {
	ID        int           `json:"id,number"`
	CreatedAt time.Time     `json:"created_at,string"`
	Amount    float64       `json:"amount,string"`
	Balance   float64       `json:"balance,string"`
	Type      string        `json:"type"`
	Details   LedgerDetails `json:"details"`
}

// LedgerDetails - If an entry is the result of a trade (match, fee), the details field will contain additional information about the trade.
type LedgerDetails struct {
	OrderID   string `json:"order_id"`
	TradeID   string `json:"trade_id"`
	ProductID string `json:"product_id"`
}

// LedgerClient API client
type LedgerClient struct {
	client *Client
}

// ListAccountLedger pulls the ledger resource for an account
func (c *LedgerClient) ListAccountLedger(id string) *Cursor {
	return NewCursor(c.client, "GET", fmt.Sprintf("/accounts/%s/ledger", id), &PaginationParams{})
}

// PullTransactionsByProduct gets ledger for an account
func (c LedgerClient) PullTransactionsByProduct(product string, t time.Time) []Ledger {
	account, err := c.client.AccountClient.FindAccountByCurrency(product)
	if err != nil {
		log.Fatal("pullTransactionByProduct:", err)
		return nil
	}

	var ledger []Ledger
	var result []Ledger
	done := false
	cursor := c.client.LedgerClient.ListAccountLedger(account.ID)
	for (cursor.HasMore) && (!done) {
		if err := cursor.NextPage(&ledger); err != nil {
			log.Fatal("pullTransactionByProduct:", err)
			return nil
		}
		for _, l := range ledger {
			if l.CreatedAt.Before(t) {
				done = true
				break
			}
			result = append(result, l)
		}
	}
	return result
}
