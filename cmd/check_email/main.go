package main

import (
	"errors"
	"flag"
	"kraken-dca-bot/internal/domain"
	"kraken-dca-bot/internal/notify"
	"log"
)

var configPath string
var newNotifier = notify.NewEmailNotifier

func init() {
	flag.StringVar(&configPath, "config", "config.yaml", "the configuration file path to test")
}

func main() {
	flag.Parse()

	config, err := domain.ParseConfig(configPath)
	if err != nil {
		log.Fatalf("failed to load configuration file : %v", err)
	}

	transaction := domain.NewTransaction("TESTPAIR")
	transaction.Fail(errors.New("test email error"))

	notifier := newNotifier(config)
	err = notifier.NotifyFailure(transaction)
	if err != nil {
		log.Fatalf("failed to notify test failure : %v", err)
	}

	log.Println("The test was run successfully")
}
