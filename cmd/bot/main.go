package main

import (
	"context"
	"flag"
	"fmt"
	krakenapi "github.com/beldur/kraken-go-api-client"
	"github.com/xhit/go-str2duration/v2"
	"kraken-dca-bot/internal/domain"
	"kraken-dca-bot/internal/kraken"
	"kraken-dca-bot/internal/notify"
	"log"
	"os"
	"time"
)

var newApi = krakenapi.New
var newTradingService = kraken.NewTrader
var newAccountService = kraken.NewAccount
var newNotifier = notify.NewEmailNotifier
var newInvestingService = kraken.NewInvestingService

var staging bool
var configPath string

func init() {
	flag.BoolVar(&staging, "staging", false, "dry run the program for testing purposes")
	flag.StringVar(&configPath, "config", "config.yaml", "the configuration file path with the DCA strategy")
}

func main() {
	flag.Parse()

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	defer func() {
		cancel()
	}()

	err := run(ctx)
	if err != nil {
		log.Printf("An error occurred : %v", err)
		os.Exit(1)
	}

	os.Exit(0)
}

func run(ctx context.Context) error {
	config, err := domain.ParseConfig(configPath)
	if err != nil {
		return fmt.Errorf("can't load the configuration : %w", err)
	}

	api := newApi(config.Kraken.Key, config.Kraken.Secret)
	tradingService := newTradingService(api, staging)
	accountService := newAccountService(api)
	notifier := newNotifier(config)
	investingService := newInvestingService(*config, accountService, tradingService, notifier)

	frequency, err := str2duration.ParseDuration(config.Frequency)
	if err != nil {
		return fmt.Errorf("cannot parse the DCA frequency environment variable : %w", err)
	}

	ticker := time.NewTicker(frequency)
	defer ticker.Stop()

	tick(investingService, notifier)

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			tick(investingService, notifier)
		}
	}
}

func tick(investingService kraken.Investor, notifier notify.Notifier) {
	transactions := investingService.Invest()

	for _, transaction := range transactions {
		if transaction.Exception != nil {
			err := notifier.NotifyFailure(transaction)
			if err != nil {
				log.Printf("An error as occurred during the failure notification : %v", err)
			}
		}
	}
}
