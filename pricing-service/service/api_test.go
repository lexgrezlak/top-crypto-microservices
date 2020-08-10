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
		name     string
		c        HttpClient
		wantByte []byte
		wantErr  error
	}{
		{
			"upstream api replies with an error",
			mockHttpClient{MockDo: func(req *http.Request) (*http.Response, error) {
				return nil, mockErr
			}},
			nil,
			mockErr,
		},
		{
			"upstream api replies with a proper response",
			mockHttpClient{MockDo: func(req *http.Request) (*http.Response, error) {
				// Marshal the mock cryptos defined before and return them
				// as a response body.

				return &http.Response{
					Body: ioutil.NopCloser(bytes.NewReader(jsonCryptoBytes)),
				}, nil
			}},
			jsonCryptoBytes,
			nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			api := NewAPI(tc.c)
			gotByte, gotErr := api.fetchCryptocurrencies()
			testUtil.Cmp(t, gotByte, tc.wantByte)
			testUtil.CmpErr(t, gotErr, tc.wantErr)
		})
	}
}
