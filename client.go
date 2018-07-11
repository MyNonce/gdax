package gdax

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

var (
	baseURL    = os.Getenv("gdax_url") // GDAX API Base URL (api.gdax.com)
	publickey  = os.Getenv("gdax_key")
	secret     = os.Getenv("gdax_secret")
	passphrase = os.Getenv("gdax_passphrase")

	// GDAXClient is the GDAX API client singleton
	GDAXClient = NewClient(baseURL, nil)
)

const (
	throttleRate time.Duration = time.Second / 3
)

// Client is used to make API call the GDAX
type Client struct {
	BaseURL    *url.URL
	httpClient *http.Client

	AccountClient *AccountClient
	LedgerClient  *LedgerClient
	FillsClient   *FillsClient
	ProductClient *ProductClient
	BookClient    *BookClient
	OrderClient   *OrderClient
}

// NewClient create a new client to make API calls
func NewClient(host string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	c := &Client{
		BaseURL: &url.URL{
			Scheme: "https",
			Host:   host,
		},
		httpClient: httpClient,
	}
	c.AccountClient = &AccountClient{client: c}
	c.LedgerClient = &LedgerClient{client: c}
	c.FillsClient = &FillsClient{client: c}
	c.ProductClient = &ProductClient{client: c}
	c.BookClient = &BookClient{client: c}
	c.OrderClient = &OrderClient{client: c}
	return c
}

// Request handles all calls to GDAX API
func (c *Client) Request(method, path string, body, result interface{}) (*http.Response, error) {
	req, err := c.newRequest(method, path, body)
	if err != nil {
		return nil, err
	}

	return c.do(req, result)
}

func (c *Client) newRequest(method, path string, params interface{}) (*http.Request, error) {
	u := fmt.Sprintf("%s%s", c.BaseURL.String(), path)
	body := bytes.NewReader(make([]byte, 0))
	var data []byte
	var err error
	if params != nil {
		data, err = json.Marshal(params)
		if err != nil {
			return nil, err
		}
		body = bytes.NewReader(data)
	}
	req, err := http.NewRequest(method, u, body)
	if err != nil {
		return nil, err
	}

	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	var signature string
	signature, err = generateSig(timestamp, path, method, string(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("CB-ACCESS-KEY", publickey)
	req.Header.Set("CB-ACCESS-SIGN", signature)
	req.Header.Set("CB-ACCESS-TIMESTAMP", timestamp)
	req.Header.Set("CB-ACCESS-PASSPHRASE", passphrase)

	return req, nil
}

func (c *Client) do(req *http.Request, result interface{}) (*http.Response, error) {

	// This will have to work for now... we are getting rate limited
	throttle := time.Tick(throttleRate)
	<-throttle

	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Fatal("client.do:", err)
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		defer resp.Body.Close()
		gdaxError := Error{}
		decoder := json.NewDecoder(resp.Body)
		if err := decoder.Decode(&gdaxError); err != nil {
			return resp, err
		}

		return resp, error(gdaxError)
	}
	if result != nil {
		decoder := json.NewDecoder(resp.Body)
		if err = decoder.Decode(result); err != nil {
			return resp, nil
		}
	}

	return resp, nil
}

func generateSig(timestamp string, path string, method string, body string) (string, error) {
	prehash := timestamp + method + path + body
	key, err := base64.StdEncoding.DecodeString(secret)
	if err != nil {
		log.Fatal("geneateSig:", err)
		return "", err
	}

	hash := hmac.New(sha256.New, key)
	_, err = hash.Write([]byte(prehash))
	if err != nil {
		log.Fatal("generateSig:", err)
		return "", err
	}

	return base64.StdEncoding.EncodeToString(hash.Sum(nil)), nil
}

func generateCurlRequest(path string, method string, req *http.Request) {
	curl := "curl "
	for k, v := range req.Header {
		curl = curl + "-H '" + k + ":" + v[0] + "' "
	}

	curl = curl + "-X " + method + " " + path
	log.Println(curl)
}
