package service

import "errors"

type mockAPI struct {
	MockFetchCryptocurrencies  func() ([]byte, error)
	ProcessCryptocurrencyBytes func() ([]Cryptocurrency, error)
}

func (api *mockAPI) FetchCryptocurrencies() ([]byte, error) {
	if api.MockFetchCryptocurrencies != nil {
		return api.MockFetchCryptocurrencies()
	}
	return nil, errors.New("something went wrong")
}
