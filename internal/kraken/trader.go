package kraken

//go:generate mockgen -destination=../mocks/mock_trading_service.go -package=mocks . Trader

import (
	"context"
	"fmt"
	"kraken-dca-bot/internal/domain"
	"log"
	"strconv"
)

type Trader interface {
	PlaceOrder(ctx context.Context, pair domain.DCAPair) error
	Fee(pair string) (float64, error)
	AskPrice(pair string) (float64, error)
}

type tradingService struct {
	api     ApiInterface
	staging bool
}

func NewTrader(api ApiInterface, staging bool) Trader {
	return &tradingService{
		api:     api,
		staging: staging,
	}
}

// Fee Get fee percentage for the given pair
func (t tradingService) Fee(pair string) (float64, error) {
	tradeVolume, err := t.api.Query("TradeVolume", map[string]string{"pair": pair, "fee-info": "true"})
	if err != nil {
		return -1, err
	}
	feePercentage, err := strconv.ParseFloat(extractData(tradeVolume, "fees", pair, "fee").(string), 64)
	if err != nil {
		return -1, err
	}

	log.Printf("[%s] Fee: %.3f%%", pair, feePercentage/100)

	return feePercentage, nil
}

// AskPrice Get the latest ticker information
func (t tradingService) AskPrice(pair string) (float64, error) {
	ticker, err := t.api.Query("Ticker", map[string]string{
		"pair": pair,
	})
	if err != nil {
		return -1, err
	}

	askPrice, err := strconv.ParseFloat(extractData(ticker, pair, "a").([]interface{})[0].(string), 64)
	if err != nil {
		return -1, err
	}

	log.Printf("[%s] Ask price : %.2f€", pair, askPrice)

	return askPrice, nil
}

// PlaceOrder Place an order for the given pair.
// The amount is specified in the DCAPair and it represent the total invested amount (token price + fees).
// The `ctx` context contains the transaction to update at the "transaction" key.
func (t tradingService) PlaceOrder(ctx context.Context, pair domain.DCAPair) error {
	transaction := ctx.Value("transaction").(*domain.Transaction)

	feePercentage, err := t.Fee(pair.Pair)
	if err != nil {
		return err
	}

	askPrice, err := t.AskPrice(pair.Pair)
	if err != nil {
		return err
	}

	orderVolume := pair.Amount / askPrice
	fee := orderVolume * feePercentage / 100
	orderVolume = orderVolume - fee
	order, err := t.api.AddOrder(pair.Pair, "buy", "market", fmt.Sprintf("%f", orderVolume), map[string]string{
		"expiretm": "+300",
		"validate": strconv.FormatBool(t.staging),
	})
	if err != nil {
		return err
	}

	orderTransactionId := "STAGED"
	if !t.staging {
		orderTransactionId = order.TransactionIds[0]
	}

	transaction.Complete(orderTransactionId, askPrice, orderVolume, fee)

	return nil
}

func extractData(data interface{}, fields ...string) interface{} {
	fieldData := data
	for _, field := range fields {
		fieldData = fieldData.(map[string]interface{})[field]
	}

	return fieldData
}
