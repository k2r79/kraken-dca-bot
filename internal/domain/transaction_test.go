package domain

import (
	"errors"
	"testing"
)

func TestTransactionComplete(t *testing.T) {
	transaction := NewTransaction("XETHZEUR")
	transaction.Complete("TXID", 123.45, 543.21, 0.123)

	if transaction.Id != "TXID" ||
		transaction.MarketPrice != 123.450000 ||
		transaction.Amount != 543.210000 ||
		transaction.Fee != 0.123000 ||
		transaction.Exception != nil {
		t.Errorf("Transaction values aren't correct %v", transaction)
	}
}

func TestTransactionFail(t *testing.T) {
	transaction := NewTransaction("XETHZEUR")
	err := errors.New("test error")
	transaction.Fail(err)

	if transaction.Exception != err {
		t.Errorf("Transaction error isn't correct %v", transaction)
	}
}

func TestTransactionString(t *testing.T) {
	transaction := NewTransaction("XETHZEUR")
	transaction.Complete("TXID", 123.45, 543.21, 0.123)

	if transaction.String() != "[TXID][XETHZEUR] 543.210000 at 123.450000 with 0.123000 fee" {
		t.Errorf("Transaction string isn't correct %v", transaction.String())
	}
}
