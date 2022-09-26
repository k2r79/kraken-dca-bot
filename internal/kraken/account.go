package kraken

import "reflect"

//go:generate mockgen -destination=../mocks/mock_account_service.go -package=mocks . Account

type Account interface {
	Balance(currency string) (float64, error)
}

type AccountService struct {
	api ApiInterface
}

func NewAccount(api ApiInterface) Account {
	return AccountService{api: api}
}

// Balance Get the Kraken account balance in Euros
func (a AccountService) Balance(currency string) (float64, error) {
	balance, err := a.api.Balance()
	if err != nil {
		return -1, err
	}

	return reflect.ValueOf(*balance).FieldByName(currency).Float(), nil
}
