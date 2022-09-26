package notify

import "kraken-dca-bot/internal/domain"

type Notifier interface {
	NotifyFailure(transaction *domain.Transaction) error
}
