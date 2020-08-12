package service

import (
	"errors"
	"net/http"
)

//type mockAPI struct {
//	MockFetchCryptocurrencies  func() ([]byte, error)
//	ProcessCryptocurrencyBytes func() ([]Cryptocurrency, error)
//}

//func (api *mockAPI) FetchCryptocurrencies() ([]byte, error) {
//	if api.MockFetchCryptocurrencies != nil {
//		return api.MockFetchCryptocurrencies()
//	}
//	return nil, errors.New("something went wrong")
//}

type mockHttpClient struct {
	MockDo func(req *http.Request) (*http.Response, error)
}

// We leave the mock function implementation to the test.
// By default it's gonna return an error
func (c mockHttpClient) Do(req *http.Request) (*http.Response, error) {
	if c.MockDo != nil {
		return c.MockDo(req)
	}
	return nil, errors.New("something went wrong")
}
