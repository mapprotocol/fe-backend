package binance

import (
	"context"
	"fmt"
	connector "github.com/binance/binance-connector-go"
	"os"
	"testing"
)

const (
	OrderSideBuy  = "BUY"
	OrderSideSELL = "SELL"
)

const (
	OrderTypeLimit  = "LIMIT"
	OrderTypeMarket = "MARKET"
)

var baseURL = "https://testnet.binance.vision"

func TestMain(m *testing.M) {
	apiKey := os.Getenv("SPOT_API_KEY")
	secretKey := os.Getenv("SPOT_SECRET_KEY")
	Init(apiKey, secretKey, baseURL)
	m.Run()
}

//func TestAsset(t *testing.T) {
//	client.Debug = true
//	order, err := client.NewUserAssetService().Asset("USDT").Do(context.Background())
//	if err != nil {
//		t.Fatal(err)
//	}
//	fmt.Println(connector.PrettyPrint(order))
//}

func TestBuy(t *testing.T) {
	client.Debug = true
	order, err := client.NewCreateOrderService().Symbol("BTCUSDT").Side(OrderSideBuy).Type(OrderTypeMarket).
		Quantity(0.001).Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(connector.PrettyPrint(order))
}

//{
//"symbol": "BTCUSDT",
//"orderId": 375527,
//"orderListId": -1,
//"clientOrderId": "QNEvwzpRmL0ZRjDhewRqRB",
//"transactTime": 1717658076505,
//"price": "0.00000000",
//"origQty": "0.00100000",
//"executedQty": "0.00100000",
//"cummulativeQuoteQty": "70.92586000",
//"status": "FILLED",
//"timeInForce": "GTC",
//"type": "MARKET",
//"side": "BUY",
//"workingTime": 1717658076505,
//"selfTradePreventionMode": "EXPIRE_MAKER",
//"fills": [
//{
//"price": "70925.86000000",
//"qty": "0.00100000",
//"commission": "0.00000000",
//"commissionAsset": "BTC",
//"tradeId": 73680
//}
//]
//}

func TestSell(t *testing.T) {
	client.Debug = true
	order, err := client.NewCreateOrderService().Symbol("BTCUSDT").Side(OrderSideSELL).Type(OrderTypeMarket).
		Quantity(0.001).Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(connector.PrettyPrint(order))
}

/*
{
	"symbol": "BTCUSDT",
	"orderId": 375937,
	"orderListId": -1,
	"clientOrderId": "pyKCOTzV4Lmzc2SepX8hfD",
	"transactTime": 1717658152146,
	"price": "0.00000000",
	"origQty": "0.00100000",
	"executedQty": "0.00100000",
	"cummulativeQuoteQty": "70.92585000",
	"status": "FILLED",
	"timeInForce": "GTC",
	"type": "MARKET",
	"side": "SELL",
	"workingTime": 1717658152146,
	"selfTradePreventionMode": "EXPIRE_MAKER",
	"fills": [
		{
			"price": "70925.85000000",
			"qty": "0.00100000",
			"commission": "0.00000000",
			"commissionAsset": "USDT",
			"tradeId": 73716
		}
	]
}
*/

func TestGetAllOrders(t *testing.T) {
	client.Debug = true
	gotOrder, err := client.NewGetAllOrdersService().Symbol("BTCUSDT").Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(connector.PrettyPrint(gotOrder))
}

/*
[
	{
		"symbol": "BTCUSDT",
		"listClientOrderId": "",
		"orderId": 375527,
		"orderListId": -1,
		"clientOrderId": "QNEvwzpRmL0ZRjDhewRqRB",
		"price": "0.00000000",
		"origQty": "0.00100000",
		"executedQty": "0.00100000",
		"cumulativeQuoteQty": "",
		"status": "FILLED",
		"timeInForce": "GTC",
		"type": "MARKET",
		"side": "BUY",
		"stopPrice": "0.00000000",
		"icebergQty": "0.00000000",
		"time": 1717658076505,
		"updateTime": 1717658076505,
		"isWorking": true,
		"origQuoteOrderQty": "0.00000000",
		"workingTime": 1717658076505,
		"selfTradePreventionMode": "EXPIRE_MAKER"
	},
	{
		"symbol": "BTCUSDT",
		"listClientOrderId": "",
		"orderId": 375937,
		"orderListId": -1,
		"clientOrderId": "pyKCOTzV4Lmzc2SepX8hfD",
		"price": "0.00000000",
		"origQty": "0.00100000",
		"executedQty": "0.00100000",
		"cumulativeQuoteQty": "",
		"status": "FILLED",
		"timeInForce": "GTC",
		"type": "MARKET",
		"side": "SELL",
		"stopPrice": "0.00000000",
		"icebergQty": "0.00000000",
		"time": 1717658152146,
		"updateTime": 1717658152146,
		"isWorking": true,
		"origQuoteOrderQty": "0.00000000",
		"workingTime": 1717658152146,
		"selfTradePreventionMode": "EXPIRE_MAKER"
	}
]
*/

func TestGetOrder(t *testing.T) {
	client.Debug = true
	gotOrder, err := client.NewGetOrderService().Symbol("BTCUSDT").OrderId(375937).Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(connector.PrettyPrint(gotOrder))
}

/*
{
	"symbol": "BTCUSDT",
	"orderId": 375937,
	"orderListId": -1,
	"clientOrderId": "pyKCOTzV4Lmzc2SepX8hfD",
	"price": "0.00000000",
	"origQty": "0.00100000",
	"executedQty": "0.00100000",
	"cumulativeQuoteQty": "",
	"status": "FILLED",
	"timeInForce": "GTC",
	"type": "MARKET",
	"side": "SELL",
	"stopPrice": "0.00000000",
	"icebergQty": "0.00000000",
	"time": 1717658152146,
	"updateTime": 1717658152146,
	"isWorking": true,
	"workingTime": 1717658152146,
	"origQuoteOrderQty": "0.00000000",
	"selfTradePreventionMode": "EXPIRE_MAKER"
}
*/

func TestGetMyTradesService(t *testing.T) {
	client.Debug = true
	gotOrder, err := client.NewGetMyTradesService().Symbol("BTCUSDT").OrderId(375937).Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(connector.PrettyPrint(gotOrder))
}

/*
[
	{
		"id": 73716,
		"symbol": "BTCUSDT",
		"orderId": 375937,
		"orderListId": -1,
		"price": "70925.85000000",
		"qty": "0.00100000",
		"quoteQty": "70.92585000",
		"commission": "0.00000000",
		"commissionAsset": "USDT",
		"time": 1717658152146,
		"isBuyer": false,
		"isMaker": false,
		"isBestMatch": true
	}
]
*/
