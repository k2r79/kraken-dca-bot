package main

import (
	"github.com/golang/mock/gomock"
	"kraken-dca-bot/internal/domain"
	"kraken-dca-bot/internal/mocks"
	"kraken-dca-bot/internal/notify"
	"testing"
)

func TestCheckEmail(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	notifier := mocks.NewMockNotifier(controller)
	newNotifier = func(config *domain.Config) notify.Notifier {
		return notifier
	}

	configPath = "../../test/data/bot-test-config.yaml"

	var transaction *domain.Transaction
	notifier.EXPECT().NotifyFailure(gomock.Any()).Do(func(tx *domain.Transaction) {
		transaction = tx
	})

	main()

	if transaction.Pair != "TESTPAIR" {
		t.Errorf("transaction pair isn't correct %v", transaction.Pair)
	}

	if transaction.Exception.Error() != "test email error" {
		t.Errorf("transaction exception isn't correct %v", transaction.Exception)
	}
}
