package domain

import (
	"fmt"
	"time"
)

type Transaction struct {
	Id          string
	Date        time.Time
	Pair        string
	MarketPrice float64
	Amount      float64
	Fee         float64
	Exception   error
}

func NewTransaction(pair string) *Transaction {
	t := &Transaction{}
	t.Date = time.Now()
	t.Pair = pair

	return t
}

func (t *Transaction) Complete(id string, marketPrice float64, amount float64, fee float64) *Transaction {
	t.Id = id
	t.MarketPrice = marketPrice
	t.Amount = amount
	t.Fee = fee

	return t
}

func (t *Transaction) Fail(exception error) *Transaction {
	t.Exception = exception

	return t
}

func (t *Transaction) String() string {
	return fmt.Sprintf("[%s][%s] %f at %f with %f fee", t.Id, t.Pair, t.Amount, t.MarketPrice, t.Fee)
}
