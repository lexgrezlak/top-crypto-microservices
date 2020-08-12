// API docs: https://coinmarketcap.com/api/documentation/v1/#section/Quick-Start-Guide

package service

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

const (
	// Normally these should be environment variables.
	UPSTREAM_API_KEY = "60c9c458-d3f0-47c7-8b60-6389c5cf9124"
	UPSTREAM_API_URL = "https://pro-api.coinmarketcap.com/v1/cryptocurrency/listings/latest"
	CURRENCY         = "USD"
	OFFSET           = "1"
	MAX_LIMIT 		 = 5000
)

// We define our own interface so that we can mock it,
// and therefore test our fetch functions.
type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// There are also other fields available. Check the docs for more information.
type Currency struct {
	Price float32 `json:"price"`
}

type Quote struct {
	USD Currency `json:"USD"`
	// We could also define fields for other currencies such as EUR
	// but USD is the one we're using - set in the `const` variable above.
}

type Cryptocurrency struct {
	Symbol string `json:"symbol"`
	Quote  Quote  `json:"quote"`
}

type fetchCryptocurrenciesBody struct {
	Cryptocurrencies []Cryptocurrency `json:"data"`
}

type api struct {
	c HttpClient
}

type CryptocurrencyDatastore interface {
	FetchCryptocurrencies(limit int) ([]byte, error)
	ProcessCryptocurrencyBytes() ([]Cryptocurrency, error)
	GetCryptocurrencies(limit int) ([]Cryptocurrency, error)
}

// We pass our custom HttpClient to enable mocking.
func NewAPI(c HttpClient) *api {
	return &api{c}
}

func (api *api) GetCryptocurrencies(limit int) ([]Cryptocurrency, error) {
	bytes, err := api.FetchCryptocurrencies(limit)
	if err != nil {
		log.Fatalf("Failed to fetch cryptocurrencies: %v", err)
		return nil, err
	}

	cryptos, err := api.ProcessCryptocurrencyBytes(bytes)
	if err != nil {
		log.Fatalf("Failed to process cryptocurrency bytes: %v", err)
		return nil, err
	}
	return cryptos, nil
}

// Unmarshals the response body received from fetch and returns the proper cryptocurrencies.
func (api *api) ProcessCryptocurrencyBytes(bytes []byte) ([]Cryptocurrency, error) {
	var body fetchCryptocurrenciesBody
	if err := json.Unmarshal(bytes, &body); err != nil {
		log.Printf("Failed to unmarshal fetch cryptocurrencies response: %v", err)
		return nil, err
	}

	return body.Cryptocurrencies, nil
}

//
// Fetches data from the CoinMarketCap api.
// Example response:
// {
//  "status": {
//    "timestamp": "2020-08-10T15:00:49.040Z",
//    "error_code": 0,
//    "error_message": null,
//    "elapsed": 100,
//    "credit_count": 1,
//    "notice": null
//  },
//  "data": [
//    {
//      "id": 1,
//      "name": "Bitcoin",
//      "symbol": "BTC",
//      "slug": "bitcoin",
//      "num_market_pairs": 8548,
//      "date_added": "2013-04-28T00:00:00.000Z",
//      "tags": [
//        "mineable",
//        "sha-256",
//        "state-channels",
//        "pow",
//        "store-of-value"
//      ],
//      "max_supply": 21000000,
//      "circulating_supply": 18456718,
//      "total_supply": 18456718,
//      "platform": null,
//      "cmc_rank": 1,
//      "last_updated": "2020-08-10T14:59:36.000Z",
//      "quote": {
//        "USD": {
//          "price": 11923.3731273,
//          "volume_24h": 25515630084.5008,
//          "percent_change_1h": 0.104577,
//          "percent_change_24h": 2.62422,
//          "percent_change_7d": 5.67848,
//          "market_cap": 220066335419.3542,
//          "last_updated": "2020-08-10T14:59:36.000Z"
//        }
//      }
//    },
//    {
//      "id": 1027,
//      "name": "Ethereum",
//      ...
func (api *api) FetchCryptocurrencies(limit int) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, UPSTREAM_API_URL, nil)
	if err != nil {
		log.Printf("Failed to create new request: %v", err)
		return nil, err
	}

	// Set up the required headers. You can get more info from the API docs.
	req.Header.Set("X-CMC_PRO_API_KEY", UPSTREAM_API_KEY)
	req.Header.Set("Accept", "application/json")

	// Set query parameters.
	q := url.Values{}
	q.Add("start", OFFSET)
	q.Add("limit", strconv.Itoa(limit))
	q.Add("convert", CURRENCY)

	req.URL.RawQuery = q.Encode()

	// Make the request.
	res, err := api.c.Do(req)
	if err != nil {
		log.Printf("Failed to fetch data from upstream API: %v", err)
		return nil, err
	}

	// Read the response body.
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("Failed to read response body: %v", err)
		return nil, err
	}

	return b, nil
}
