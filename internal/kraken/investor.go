package kraken

//go:generate mockgen -destination=../mocks/mock_investing_service.go -package=mocks . Investor

import (
	"context"
	"errors"
	"fmt"
	"kraken-dca-bot/internal/domain"
	"kraken-dca-bot/internal/notify"
	"log"
	"time"
)

type Investor interface {
	Invest() []*domain.Transaction
}

type investingService struct {
	config         domain.Config
	accountService Account
	tradingService Trader
	notifier       notify.Notifier
}

func NewInvestingService(config domain.Config, accountService Account, tradingService Trader, notifier notify.Notifier) Investor {
	return investingService{
		config:         config,
		accountService: accountService,
		tradingService: tradingService,
		notifier:       notifier,
	}
}

func (i investingService) Invest() []*domain.Transaction {
	start := time.Now()
	transactions := make([]*domain.Transaction, len(i.config.Pairs))

	for index, pair := range i.config.Pairs {
		log.Printf("Trading %s...", pair.Pair)

		transactions[index] = i.investInPair(pair)
	}

	log.Printf("Execution time : %s", time.Since(start))

	return transactions
}

func (i investingService) investInPair(pair domain.DCAPair) *domain.Transaction {
	transaction := domain.NewTransaction(pair.Pair)
	ctx := context.Background()
	ctx = context.WithValue(ctx, "transaction", transaction)

	var err error
	defer func() {
		if err != nil {
			transaction.Fail(err)
		}
	}()

	accountBalance, err := i.accountService.Balance(i.config.Currency)
	if err != nil {
		err = fmt.Errorf("account balance cannot be collected : %w", err)

		return transaction
	}
	log.Printf("Account balance : %.2f€", accountBalance)

	if accountBalance < pair.Amount {
		err = errors.New("account balance is less than the pair DCA amount, stopping the process")

		return transaction
	}

	err = i.tradingService.PlaceOrder(ctx, pair)
	log.Println(transaction)
	if err != nil {
		notifyErr := i.notifier.NotifyFailure(transaction)
		if notifyErr != nil {
			err = fmt.Errorf("failed to notify %s transaction failure : %w", pair.Pair, notifyErr)
		}
		err = fmt.Errorf("could not place order on %s : %w", pair.Pair, err)
	}

	return transaction
}
