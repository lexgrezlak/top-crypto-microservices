// API docs: https://min-api.cryptocompare.com/documentation

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
	// Normally it should be an environment variable.
	UPSTREAM_API_URL = "https://min-api.cryptocompare.com/data/top/mktcapfull"
	// Currency is the `tsym` parameter for CryptoCompare.
	CURRENCY         = "USD"
	// 100 is the max limit CryptoCompare allows to fetch for one page.
	PAGE_SIZE 	      = 100
	// The number of results to fetch set by the problem specification.
	//
	// Likely it would've been better to pass the fetch size through the consumer,
	// however, it's our first shot with RabbitMQ, and we currently don't know
	// how to do this.
	FETCH_SIZE = 200
)

// We define our own interface so that we can mock it,
// and therefore test our fetch functions.
type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// We pass our custom HttpClient to enable mocking.
func NewAPI(c HttpClient) *api {
	return &api{c}
}

type Cryptocurrency struct {
	// We define only the properties we're interested in, that is:
	// we only need the names of cryptocurrencies for this problem.
	CoinInfo struct{
		Name string `json:"Name"`
	} `json:"CoinInfo"`
}

type fetchCryptocurrenciesBody struct {
	Data []Cryptocurrency `json:"Data"`
}

type api struct {
	c HttpClient
}

type CryptocurrencyDatastore interface {
	FetchCryptocurrencies() ([]byte, error)
	GetSymbolsFromCryptocurrencyBytes() ([]string, error)
	GetCryptocurrencySymbols() ([]string, error)
}

func (api *api) GetCryptocurrencySymbols() ([]string, error) {
	var symbols []string
	pagesToFetchCount := FETCH_SIZE / PAGE_SIZE

	// Fetch the cryptos, process them, and append to the symbols array, for the number of pages
	for pageNum := 0; pageNum < pagesToFetchCount; pageNum++ {
		bytes, err := api.FetchCryptocurrencies(pageNum)
		if err != nil {
			log.Printf("Failed to fetch cryptocurrencies: %v", err)
			return nil, err
		}
		currSymbols, err := api.GetSymbolsFromCryptocurrencyBytes(bytes)
		if err != nil {
			log.Printf("Failed to get symbols from cryptocurrency bytes: %v", err)
			return nil, err
		}
		symbols = append(currSymbols, symbols...)
	}

	return symbols, nil
}

// Unmarshals the response body received from fetch and returns the proper cryptocurrencies.
func (api *api) GetSymbolsFromCryptocurrencyBytes(bytes []byte) ([]string, error) {
	var body fetchCryptocurrenciesBody
	if err := json.Unmarshal(bytes, &body); err != nil {
		log.Printf("Failed to unmarshal fetch cryptocurrencies response: %v", err)
		return nil, err
	}

	var names []string

	for _, cryptocurrency := range body.Data {
		names = append(names, cryptocurrency.CoinInfo.Name)
	}

	log.Println(names)

	return names, nil
}

// Fetches data from the CryptoCompare api.
// Example response:
// {
//  "Message": "Success",
//  "Type": 100,
//  "SponsoredData": [],
//  "Data": [
//    {
//      "CoinInfo": {
//        "Id": "1182",
//        "Name": "BTC",
//        "FullName": "Bitcoin",
//        "Internal": "BTC",
// 		  ...
// 	  	},
//      "RAW": {
//        "USD": {
//          "TYPE": "5",
//          "MARKET": "CCCAGG",
//          "FROMSYMBOL": "BTC",
//          "TOSYMBOL": "USD",
//          "FLAGS": "2052",
//          "PRICE": 11511.4,
//			...
// 		 },
//      "DISPLAY": {
//        "USD": {
//          "FROMSYMBOL": "Ƀ",
//          "TOSYMBOL": "$",
//          "MARKET": "CryptoCompare Index",
//          "PRICE": "$ 11,511.4",
//			...
//   	  }
//    },
//    {
//      "CoinInfo": {
//        "Id": "934865",
//        "Name": "TNCC",
//        "FullName": "TNC Coin",
//		  ...
func (api *api) FetchCryptocurrencies(page int) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, UPSTREAM_API_URL, nil)
	if err != nil {
		log.Printf("Failed to create new request: %v", err)
		return nil, err
	}

	// Set query parameters.
	q := url.Values{}
	q.Add("page", strconv.Itoa(page))
	q.Add("limit", strconv.Itoa(PAGE_SIZE))
	// tsym is the parameter that specifies currency received from the API.
	q.Add("tsym", CURRENCY)

	// Encode the query parameters.
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
