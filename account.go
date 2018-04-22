package gdax

import (
	"fmt"
	"log"
	"strings"
)

// Account represent the Gdax account
type Account struct {
	ID        string  `json:"id"`
	Currency  string  `json:"currency"`
	Balance   float64 `json:"balance,string"`
	Available float64 `json:"available,string"`
	Hold      float64 `json:"hold,string"`
	ProfileID string  `json:"profile_id"`
}

// AccountClient API client
type AccountClient struct {
	client *Client
}

// GetAccounts fetches accounts
func (ac AccountClient) GetAccounts() ([]Account, error) {
	var accounts []Account
	_, err := ac.client.Request("GET", "/accounts", nil, &accounts)
	return accounts, err
}

// GetAccount returns a single Account by its ID
func (ac AccountClient) GetAccount(id string) (Account, error) {
	var account Account

	resource := fmt.Sprintf("/accounts/%s", id)
	_, err := ac.client.Request("GET", resource, nil, account)
	return account, err
}

// FindAccountByCurrency fetches the account for a particular currency, ETH, BTC, LTC, USD
func (ac AccountClient) FindAccountByCurrency(currency string) (*Account, error) {
	accounts, err := ac.GetAccounts()
	if err != nil {
		log.Fatal("findAccountByCurrency:", err)
		return nil, err
	}
	for _, acc := range accounts {
		if strings.Compare(acc.Currency, currency) == 0 {
			return &acc, nil
		}
	}
	return nil, nil
}
