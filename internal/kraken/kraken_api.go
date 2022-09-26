package kraken

//go:generate mockgen -destination=../mocks/mock_kraken_api.go -package=mocks . ApiInterface

import krakenapi "github.com/beldur/kraken-go-api-client"

type ApiInterface interface {
	Balance() (*krakenapi.BalanceResponse, error)
	Query(method string, data map[string]string) (interface{}, error)
	AddOrder(pair string, direction string, orderType string, volume string, args map[string]string) (*krakenapi.AddOrderResponse, error)
}
