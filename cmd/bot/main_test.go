package main

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"kraken-dca-bot/internal/domain"
	"kraken-dca-bot/internal/kraken"
	"kraken-dca-bot/internal/mocks"
	"kraken-dca-bot/internal/notify"
	"strings"
	"testing"
	"time"
)

var tradingService *mocks.MockTrader
var accountService *mocks.MockAccount
var notifier *mocks.MockNotifier
var investingService *mocks.MockInvestor

func setup(t *testing.T) func() {
	controller := gomock.NewController(t)

	staging = true
	configPath = "../../test/data/bot-test-config.yaml"

	tradingService = mocks.NewMockTrader(controller)
	newTradingService = func(api kraken.ApiInterface, staging bool) kraken.Trader {
		return tradingService
	}

	accountService = mocks.NewMockAccount(controller)
	newAccountService = func(api kraken.ApiInterface) kraken.Account {
		return accountService
	}

	notifier = mocks.NewMockNotifier(controller)
	newNotifier = func(config *domain.Config) notify.Notifier {
		return notifier
	}

	investingService = mocks.NewMockInvestor(controller)
	newInvestingService = func(config domain.Config, accountService kraken.Account, tradingService kraken.Trader, notifier notify.Notifier) kraken.Investor {
		return investingService
	}

	return controller.Finish
}

func TestBotSuccess(t *testing.T) {
	cleanUp := setup(t)
	defer cleanUp()

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	transactions := []*domain.Transaction{
		{
			Id:        "TXID1",
			Date:      time.Time{},
			Exception: nil,
		},
		{
			Id:        "TXID2",
			Date:      time.Time{},
			Exception: nil,
		},
	}

	investingService.EXPECT().Invest().Return(transactions)
	investingService.EXPECT().Invest().DoAndReturn(func() []*domain.Transaction {
		cancel()
		return transactions
	})

	err := run(ctx)
	if err != nil {
		t.Errorf("An unexpected error occurred : %v", err)
	}
}

func TestBotInvestFail(t *testing.T) {
	cleanUp := setup(t)
	defer cleanUp()

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	transactions := []*domain.Transaction{
		{
			Id:        "TXID1",
			Date:      time.Time{},
			Exception: nil,
		},
		{
			Id:        "TXID2",
			Date:      time.Time{},
			Exception: errors.New("transaction error"),
		},
	}

	investingService.EXPECT().Invest().DoAndReturn(func() []*domain.Transaction {
		cancel()
		return transactions
	})
	notifier.EXPECT().NotifyFailure(transactions[1]).Return(nil)

	err := run(ctx)
	if err != nil {
		t.Errorf("An unexpected error occurred : %v", err)
	}
}

func TestBotInvestFailNotifierFail(t *testing.T) {
	cleanUp := setup(t)
	defer cleanUp()

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	transactions := []*domain.Transaction{
		{
			Id:        "TXID1",
			Date:      time.Time{},
			Exception: errors.New("transaction error"),
		},
		{
			Id:        "TXID2",
			Date:      time.Time{},
			Exception: nil,
		},
	}

	investingService.EXPECT().Invest().DoAndReturn(func() []*domain.Transaction {
		cancel()
		return transactions
	})
	notifier.EXPECT().NotifyFailure(transactions[0]).Return(errors.New("notify error"))
	notifier.EXPECT().NotifyFailure(transactions[1]).Return(nil).Times(0)

	err := run(ctx)
	if err != nil {
		t.Errorf("An unexpected error occurred : %v", err)
	}
}

func TestBotParseConfigFail(t *testing.T) {
	staging = true
	configPath = "/missing-config.yaml"

	err := run(context.Background())
	if err == nil || !strings.HasPrefix(err.Error(), "can't load the configuration :") {
		t.Errorf("An unexpected error was raised : %v", err)
	}
}

func TestBotFrequencyParseFail(t *testing.T) {
	cleanUp := setup(t)
	defer cleanUp()

	staging = true
	configPath = "../../test/data/invalid-frequency.yaml"

	err := run(context.Background())
	if err == nil || !strings.HasPrefix(err.Error(), "cannot parse the DCA frequency environment variable :") {
		t.Errorf("An unexpected error was raised : %v", err)
	}
}
