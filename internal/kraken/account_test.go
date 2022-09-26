package kraken

import (
	"errors"
	krakenapi "github.com/beldur/kraken-go-api-client"
	"github.com/golang/mock/gomock"
	"kraken-dca-bot/internal/mocks"
	"testing"
)

func TestBalanceSuccess(t *testing.T) {
	cases := []struct {
		error   error
		balance float64
	}{
		{nil, 10.64},
		{errors.New("balance error"), -1},
	}

	controller := gomock.NewController(t)
	defer controller.Finish()

	for _, c := range cases {
		krakenApi := mocks.NewMockApiInterface(controller)
		accountService := NewAccount(krakenApi)

		krakenApi.EXPECT().Balance().Return(&krakenapi.BalanceResponse{ZEUR: c.balance}, c.error)

		balance, err := accountService.Balance("ZEUR")

		if err != c.error {
			t.Errorf("An unexpected error was returned : %v", err)
		}

		if balance != c.balance {
			t.Errorf("Balance returned %f instead of 10.64", balance)
		}
	}
}
