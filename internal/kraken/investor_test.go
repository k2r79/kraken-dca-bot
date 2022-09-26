package kraken

import (
	"errors"
	"github.com/golang/mock/gomock"
	"kraken-dca-bot/internal/domain"
	"kraken-dca-bot/internal/mocks"
	"testing"
)

var config = domain.Config{
	Currency: "ZEUR",
	Pairs: []domain.DCAPair{
		{
			Pair:   "XETHZEUR",
			Amount: 20.00,
		},
		{
			Pair:   "XXBTZEUR",
			Amount: 10.00,
		},
	},
}

func TestInvestSuccess(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	accountService := mocks.NewMockAccount(controller)
	tradingService := mocks.NewMockTrader(controller)
	notifier := mocks.NewMockNotifier(controller)

	investingService := NewInvestingService(config, accountService, tradingService, notifier)

	accountService.EXPECT().Balance("ZEUR").Return(35.23, nil)
	accountService.EXPECT().Balance("ZEUR").Return(15.23, nil)
	tradingService.EXPECT().PlaceOrder(gomock.Any(), config.Pairs[0]).Return(nil)
	tradingService.EXPECT().PlaceOrder(gomock.Any(), config.Pairs[1]).Return(nil)
	notifier.EXPECT().NotifyFailure(gomock.Any()).Times(0)

	transactions := investingService.Invest()

	if len(transactions) != 2 {
		t.Errorf("Transaction count is wrong : %v", len(transactions))
	}

	if transactions[0].Pair != "XETHZEUR" {
		t.Errorf("First transaction pair is wrong : %v", transactions[0].Pair)
	}

	if transactions[1].Pair != "XXBTZEUR" {
		t.Errorf("Second transaction pair is wrong : %v", transactions[1].Pair)
	}
}

func TestInvestFail_BalanceFail(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	accountService := mocks.NewMockAccount(controller)
	tradingService := mocks.NewMockTrader(controller)
	notifier := mocks.NewMockNotifier(controller)

	investingService := NewInvestingService(config, accountService, tradingService, notifier)

	accountService.EXPECT().Balance("ZEUR").Return(0.0, errors.New("balance error")).Times(2)

	transactions := investingService.Invest()

	if transactions[0].Pair != "XETHZEUR" {
		t.Errorf("First transaction pair is %v", transactions[0].Pair)
	}

	if transactions[0].Exception.Error() != "account balance cannot be collected : balance error" {
		t.Errorf("First transaction exception is %v", transactions[0].Exception)
	}

	if transactions[1].Pair != "XXBTZEUR" {
		t.Errorf("Second transaction pair is %v", transactions[1].Pair)
	}

	if transactions[1].Exception.Error() != "account balance cannot be collected : balance error" {
		t.Errorf("Second transaction pair is %v", transactions[1].Exception)
	}
}

func TestInvestFail_BalanceInsufficient(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	accountService := mocks.NewMockAccount(controller)
	tradingService := mocks.NewMockTrader(controller)
	notifier := mocks.NewMockNotifier(controller)

	investingService := NewInvestingService(config, accountService, tradingService, notifier)

	accountService.EXPECT().Balance("ZEUR").Return(23.08, nil)
	tradingService.EXPECT().PlaceOrder(gomock.Any(), gomock.Any()).Return(nil)
	accountService.EXPECT().Balance("ZEUR").Return(3.08, nil)

	transactions := investingService.Invest()

	if transactions[0].Pair != "XETHZEUR" {
		t.Errorf("First transaction pair is %v", transactions[0].Pair)
	}

	if transactions[0].Exception != nil {
		t.Errorf("First transaction exception is %v", transactions[0].Exception)
	}

	if transactions[1].Pair != "XXBTZEUR" {
		t.Errorf("Second transaction pair is %v", transactions[1].Pair)
	}

	if transactions[1].Exception.Error() != "account balance is less than the pair DCA amount, stopping the process" {
		t.Errorf("Second transaction pair is %v", transactions[1].Exception)
	}
}

func TestInvestFail_PlaceOrderFail(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	accountService := mocks.NewMockAccount(controller)
	tradingService := mocks.NewMockTrader(controller)
	notifier := mocks.NewMockNotifier(controller)

	investingService := NewInvestingService(config, accountService, tradingService, notifier)

	accountService.EXPECT().Balance("ZEUR").Return(45.44, nil)
	tradingService.EXPECT().PlaceOrder(gomock.Any(), config.Pairs[0]).Return(nil)
	accountService.EXPECT().Balance("ZEUR").Return(25.44, nil)
	tradingService.EXPECT().PlaceOrder(gomock.Any(), config.Pairs[1]).Return(errors.New("place order error"))
	notifier.EXPECT().NotifyFailure(gomock.Any()).Do(func(transaction *domain.Transaction) {
		if transaction.Pair != "XXBTZEUR" {
			t.Errorf("The notifier transaction pair is %v", transaction.Pair)
		}
	}).Return(nil)

	transactions := investingService.Invest()

	if transactions[0].Pair != "XETHZEUR" {
		t.Errorf("First transaction pair is %v", transactions[0].Pair)
	}

	if transactions[0].Exception != nil {
		t.Errorf("First transaction exception is %v", transactions[0].Exception)
	}

	if transactions[1].Pair != "XXBTZEUR" {
		t.Errorf("Second transaction pair is %v", transactions[1].Pair)
	}

	if transactions[1].Exception.Error() != "could not place order on XXBTZEUR : place order error" {
		t.Errorf("Second transaction pair is %v", transactions[1].Exception)
	}
}

func TestInvestFail_PlaceOrderFailNotifierFail(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	accountService := mocks.NewMockAccount(controller)
	tradingService := mocks.NewMockTrader(controller)
	notifier := mocks.NewMockNotifier(controller)

	investingService := NewInvestingService(config, accountService, tradingService, notifier)

	accountService.EXPECT().Balance("ZEUR").Return(45.44, nil)
	tradingService.EXPECT().PlaceOrder(gomock.Any(), config.Pairs[0]).Return(nil)
	accountService.EXPECT().Balance("ZEUR").Return(25.44, nil)
	tradingService.EXPECT().PlaceOrder(gomock.Any(), config.Pairs[1]).Return(errors.New("place order error"))
	notifier.EXPECT().NotifyFailure(gomock.Any()).Return(errors.New("notifier error"))

	transactions := investingService.Invest()

	if transactions[0].Pair != "XETHZEUR" {
		t.Errorf("First transaction pair is %v", transactions[0].Pair)
	}

	if transactions[0].Exception != nil {
		t.Errorf("First transaction exception is %v", transactions[0].Exception)
	}

	if transactions[1].Pair != "XXBTZEUR" {
		t.Errorf("Second transaction pair is %v", transactions[1].Pair)
	}

	if transactions[1].Exception.Error() != "could not place order on XXBTZEUR : failed to notify XXBTZEUR transaction failure : notifier error" {
		t.Errorf("Second transaction pair is %v", transactions[1].Exception)
	}
}
