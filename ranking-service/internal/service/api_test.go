package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"testing"
	"top-coins/pricing-service/internal/testUtil"
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

	type want struct {
		bytes []byte
		err   error
	}

	testCases := []struct {
		name string
		c    HttpClient
		want want
	}{
		{
			"upstream api replies with an error",
			mockHttpClient{MockDo: func(req *http.Request) (*http.Response, error) {
				return nil, mockErr
			}},
			want{
				nil,
				mockErr,
			},
		},
		{
			"upstream api replies with a proper response",
			mockHttpClient{MockDo: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					Body: ioutil.NopCloser(bytes.NewReader(jsonCryptoBytes)),
				}, nil
			}},
			want{
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

	type want struct {
		Cryptos []Cryptocurrency
		err     error
	}

	testCases := []struct {
		name string
		want want
	}{
		{
			"valid data",
			want{
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
