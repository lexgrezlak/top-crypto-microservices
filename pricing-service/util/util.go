package util

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

const (
	// Normally it should be an environment variable.
	UPSTREAM_API_KEY = "60c9c458-d3f0-47c7-8b60-6389c5cf9124"
	UPSTREAM_API_URL = "https://pro-api.coinmarketcap.com/v1/cryptocurrency/listings/latest"
)

// API docs: https://coinmarketcap.com/api/documentation/v1/#section/Quick-Start-Guide
func FetchPrices() ([]byte,error) {
	c := &http.Client{
		Timeout:       15 * time.Second,
	}
	req, err := http.NewRequest(http.MethodGet, UPSTREAM_API_URL, nil)
	if err != nil {
		log.Printf("Failed to create new request: %v", err)
		return nil, err
	}

	req.Header.Set("X-CMC_PRO_API_KEY", UPSTREAM_API_KEY)
	req.Header.Set("Accept", "application/json")

	q := url.Values{}
	q.Add("start", "1")
	q.Add("limit", "5000")
	q.Add("convert", "USD")

	res, err := c.Do(req)
	if err != nil {
		log.Printf("Failed to fetch data from upstream API: %v", err)
		return nil, err
	}
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("Failed to read response body: %v", err)
		return nil, err
	}
	return b, nil
}