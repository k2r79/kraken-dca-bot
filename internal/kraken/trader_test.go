package kraken

import (
	"context"
	"errors"
	krakenapi "github.com/beldur/kraken-go-api-client"
	"github.com/golang/mock/gomock"
	"kraken-dca-bot/internal/domain"
	"kraken-dca-bot/internal/mocks"
	"strconv"
	"testing"
)

var krakenApi *mocks.MockApiInterface
var service Trader

func setup(t *testing.T) func() {
	controller := gomock.NewController(t)

	krakenApi = mocks.NewMockApiInterface(controller)
	service = NewTrader(krakenApi, false)

	return controller.Finish
}

// Fee method tests

func TestFeeSuccess(t *testing.T) {
	cleanUp := setup(t)
	defer cleanUp()

	krakenApi.EXPECT().Query(
		"TradeVolume",
		map[string]string{"pair": "TESTPAIR", "fee-info": "true"},
	).Return(
		map[string]interface{}{
			"fees": map[string]interface{}{
				"TESTPAIR": map[string]interface{}{
					"fee": "0.234",
				},
			},
		},
		nil)

	feePercentage, err := service.Fee("TESTPAIR")
	if err != nil {
		t.Errorf("An unexpected fee error occured : %v", err)
	}

	if feePercentage != 0.234 {
		t.Errorf("Fee percentage is equal to %f instead of 0.234", feePercentage)
	}
}

func TestFeeFail(t *testing.T) {
	cleanUp := setup(t)
	defer cleanUp()

	krakenApi.EXPECT().Query(
		"TradeVolume",
		map[string]string{"pair": "TESTPAIR", "fee-info": "true"},
	).Return(nil, errors.New("fee error"))

	feePercentage, err := service.Fee("TESTPAIR")
	if err == nil || err.Error() != "fee error" {
		t.Error("No error was raised by the fee API call")
	}

	if feePercentage != -1 {
		t.Error("Fee percentage isn't nil")
	}
}

func TestFeeParseFail(t *testing.T) {
	cleanUp := setup(t)
	defer cleanUp()

	krakenApi.EXPECT().Query(
		"TradeVolume",
		map[string]string{"pair": "TESTPAIR", "fee-info": "true"},
	).Return(
		map[string]interface{}{
			"fees": map[string]interface{}{
				"TESTPAIR": map[string]interface{}{
					"fee": "abc",
				},
			},
		},
		nil)

	feePercentage, err := service.Fee("TESTPAIR")
	_, ok := err.(*strconv.NumError)
	if err == nil || !ok {
		t.Error("No relevant error was raised by the fee API call")
	}

	if feePercentage != -1 {
		t.Error("Fee percentage isn't nil")
	}
}

// askPrice method tests

func TestAskPriceSuccess(t *testing.T) {
	cleanUp := setup(t)
	defer cleanUp()

	krakenApi.EXPECT().Query(
		"Ticker",
		map[string]string{"pair": "TESTPAIR"},
	).Return(
		map[string]interface{}{
			"TESTPAIR": map[string]interface{}{
				"a": []interface{}{"1545.89"},
			},
		},
		nil)

	askPrice, err := service.AskPrice("TESTPAIR")
	if err != nil {
		t.Errorf("An unexpected ask price error occured : %v", err)
	}

	if askPrice != 1545.89 {
		t.Errorf("Ask price is equal to %f instead of 1545.89", askPrice)
	}
}

func TestAskPriceFail(t *testing.T) {
	cleanUp := setup(t)
	defer cleanUp()

	krakenApi.EXPECT().Query(
		"Ticker",
		map[string]string{"pair": "TESTPAIR"},
	).Return(nil, errors.New("ask price error"))

	askPrice, err := service.AskPrice("TESTPAIR")
	if err == nil || err.Error() != "ask price error" {
		t.Errorf("No relevant ask price error occured : %v", err)
	}

	if askPrice != -1 {
		t.Errorf("Ask price is equal to %f instead of -1", askPrice)
	}
}

func TestAskPriceParseFail(t *testing.T) {
	cleanUp := setup(t)
	defer cleanUp()

	krakenApi.EXPECT().Query(
		"Ticker",
		map[string]string{"pair": "TESTPAIR"},
	).Return(
		map[string]interface{}{
			"TESTPAIR": map[string]interface{}{
				"a": []interface{}{"abc"},
			},
		},
		nil)

	askPrice, err := service.AskPrice("TESTPAIR")
	_, ok := err.(*strconv.NumError)
	if err == nil || !ok {
		t.Errorf("No relevant ask price error occured : %v", err)
	}

	if askPrice != -1 {
		t.Errorf("Ask price is equal to %f instead of -1", askPrice)
	}
}

// PlaceOrder method tests
func TestPlaceOrderSuccess(t *testing.T) {
	cleanUp := setup(t)
	defer cleanUp()

	krakenApi.EXPECT().Query(
		"TradeVolume",
		map[string]string{"pair": "TESTPAIR", "fee-info": "true"},
	).Return(
		map[string]interface{}{
			"fees": map[string]interface{}{
				"TESTPAIR": map[string]interface{}{
					"fee": "0.234",
				},
			},
		},
		nil)

	krakenApi.EXPECT().Query(
		"Ticker",
		map[string]string{"pair": "TESTPAIR"},
	).Return(
		map[string]interface{}{
			"TESTPAIR": map[string]interface{}{
				"a": []interface{}{"1545.89"},
			},
		},
		nil)

	orderResponse := &krakenapi.AddOrderResponse{
		TransactionIds: []string{"ID"},
		Description: krakenapi.OrderDescription{
			Order:        "Order description",
			PrimaryPrice: "1545.89",
		},
	}
	krakenApi.EXPECT().AddOrder(
		"TESTPAIR",
		"buy",
		"market",
		"0.012907",
		map[string]string{
			"expiretm": "+300",
			"validate": "false",
		},
	).Return(orderResponse, nil)

	transaction := domain.Transaction{}
	ctx := context.Background()
	ctx = context.WithValue(ctx, "transaction", &transaction)
	err := service.PlaceOrder(ctx, domain.DCAPair{
		Pair:   "TESTPAIR",
		Amount: 20.00,
	})

	if transaction.Id != "ID" {
		t.Errorf("The transaction ID is %v", transaction.Id)
	}

	if transaction.MarketPrice != 1545.89 {
		t.Errorf("The transaction market price is %v", transaction.MarketPrice)
	}

	if transaction.Amount != 0.012907257308087897 {
		t.Errorf("The transaction amount is %v", transaction.Amount)
	}

	if err != nil {
		t.Errorf("An unexpected error has been raised : %v", err)
	}
}

func TestPlaceOrderFeeFail(t *testing.T) {
	cleanUp := setup(t)
	defer cleanUp()

	krakenApi.EXPECT().Query(
		"TradeVolume",
		map[string]string{"pair": "TESTPAIR", "fee-info": "true"},
	).Return(nil, errors.New("fee error"))

	ctx := context.Background()
	ctx = context.WithValue(ctx, "transaction", &domain.Transaction{})
	err := service.PlaceOrder(ctx, domain.DCAPair{
		Pair:   "TESTPAIR",
		Amount: 20.00,
	})

	if err == nil || err.Error() != "fee error" {
		t.Errorf("An unexpected error has been raised : %v", err)
	}
}

func TestPlaceOrderTickerFail(t *testing.T) {
	cleanUp := setup(t)
	defer cleanUp()

	krakenApi.EXPECT().Query(
		"TradeVolume",
		map[string]string{"pair": "TESTPAIR", "fee-info": "true"},
	).Return(
		map[string]interface{}{
			"fees": map[string]interface{}{
				"TESTPAIR": map[string]interface{}{
					"fee": "0.234",
				},
			},
		},
		nil)

	krakenApi.EXPECT().Query(
		"Ticker",
		map[string]string{"pair": "TESTPAIR"},
	).Return(nil, errors.New("ticker error"))

	ctx := context.Background()
	ctx = context.WithValue(ctx, "transaction", &domain.Transaction{})
	err := service.PlaceOrder(ctx, domain.DCAPair{
		Pair:   "TESTPAIR",
		Amount: 20.00,
	})

	if err == nil || err.Error() != "ticker error" {
		t.Errorf("An unexpected error has been raised : %v", err)
	}
}

func TestPlaceOrderFail(t *testing.T) {
	cleanUp := setup(t)
	defer cleanUp()

	krakenApi.EXPECT().Query(
		"TradeVolume",
		map[string]string{"pair": "TESTPAIR", "fee-info": "true"},
	).Return(
		map[string]interface{}{
			"fees": map[string]interface{}{
				"TESTPAIR": map[string]interface{}{
					"fee": "0.234",
				},
			},
		},
		nil)

	krakenApi.EXPECT().Query(
		"Ticker",
		map[string]string{"pair": "TESTPAIR"},
	).Return(
		map[string]interface{}{
			"TESTPAIR": map[string]interface{}{
				"a": []interface{}{"1545.89"},
			},
		},
		nil)

	krakenApi.EXPECT().AddOrder(
		"TESTPAIR",
		"buy",
		"market",
		"0.012907",
		map[string]string{
			"expiretm": "+300",
			"validate": "false",
		},
	).Return(nil, errors.New("place order error"))

	ctx := context.Background()
	ctx = context.WithValue(ctx, "transaction", &domain.Transaction{})
	err := service.PlaceOrder(ctx, domain.DCAPair{
		Pair:   "TESTPAIR",
		Amount: 20.00,
	})

	if err == nil || err.Error() != "place order error" {
		t.Errorf("An unexpected error has been raised : %v", err)
	}
}
