package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"testing"
	"top-coins/pricing-service/testUtil"
)

func TestFetchCryptocurrencies(t *testing.T) {
	cryptos := []Cryptocurrency{
		{
			Symbol: "BTC",
			Quote: Quote{
				USD: Currency{
					Price: 39.21,
				},
			},
		},
		{
			Symbol: "ETH",
			Quote: Quote{
				USD: Currency{
					Price: 99.02,
				},
			},
		},
	}
	jsonCryptoBytes, err := json.Marshal(cryptos)
	if err != nil {
		t.Fatalf("Failed to marshal cryptocurrencies: %v", err)
	}
	// Mock error must be defined here, otherwise the pointers will be different
	// and CmpErr will make the test fail.
	mockErr := errors.New("hello world")

	testCases := []struct {
		name string
		c    HttpClient
		want struct {
			bytes []byte
			err   error
		}
	}{
		{
			"upstream api replies with an error",
			testUtil.MockHttpClient{MockDo: func(req *http.Request) (*http.Response, error) {
				return nil, mockErr
			}},
			struct {
				bytes []byte
				err   error
			}{
				bytes: nil,
				err:   mockErr},
		},
		{
			"upstream api replies with a proper response",
			testUtil.MockHttpClient{MockDo: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					Body: ioutil.NopCloser(bytes.NewReader(jsonCryptoBytes)),
				}, nil
			}},
			struct {
				bytes []byte
				err   error
			}{
				bytes: jsonCryptoBytes,
				err:   nil,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			api := NewAPI(tc.c)
			gotBytes, gotErr := api.FetchCryptocurrencies()
			testUtil.Cmp(t, gotBytes, tc.want.bytes)
			testUtil.CmpErr(t, gotErr, tc.want.err)
		})
	}
}

func TestGetCryptocurrencies(t *testing.T) {
	cryptos := []Cryptocurrency{
		{
			Symbol: "BTC",
			Quote: Quote{
				USD: Currency{
					Price: 39.21,
				},
			},
		},
		{
			Symbol: "ETH",
			Quote: Quote{
				USD: Currency{
					Price: 99.02,
				},
			},
		},
	}
	mockFetchRes := fetchCryptocurrenciesBody{Cryptocurrencies: cryptos}

	jsonMockFetchRes, err := json.Marshal(mockFetchRes)
	if err != nil {
		t.Fatalf("Failed to marshal cryptocurrencies: %v", err)
	}

	testCases := []struct {
		name string
		want struct {
			Cryptos []Cryptocurrency
			err     error
		}
	}{
		{
			"valid data",
			struct {
				Cryptos []Cryptocurrency
				err     error
			}{
				Cryptos: cryptos,
				err:     nil,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			api := NewAPI(nil)
			gotCryptos, gotErr := api.ProcessCryptocurrencyBytes(jsonMockFetchRes)
			testUtil.Cmp(t, gotCryptos, tc.want.Cryptos)
			testUtil.CmpErr(t, gotErr, tc.want.err)
		})
	}
}
