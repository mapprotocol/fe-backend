package logic

// ref:
// https://binance-docs.github.io/apidocs/spot/cn/#api-3
// https://github.com/binance/binance-spot-api-docs/blob/master/rest-api.md#enum-definitions
const (
	BinanceOrderSideBuy  = "BUY"
	BinanceOrderSideSELL = "SELL"
)

const (
	BinanceOrderTypeLimit  = "LIMIT"
	BinanceOrderTypeMarket = "MARKET"
)

const (
	BinanceOrderRespTypeAck    = "ACK"
	BinanceOrderRespTypeResult = "RESULT"
	BinanceOrderRespTypeFull   = "FULL"
)

const (
	BinanceOrderStatusNew             = "NEW"
	BinanceOrderStatusPartiallyFilled = "PARTIALLY_FILLED"
	BinanceOrderStatusFilled          = "FILLED"
	BinanceOrderStatusCanceled        = "CANCELED"
	BinanceOrderStatusPendingCancel   = "PENDING_CANCEL"
	BinanceOrderStatusRejected        = "REJECTED"
	BinanceOrderStatusExpired         = "EXPIRED"
	BinanceOrderStatusExpiredInMatch  = "EXPIRED_IN_MATCH"
)
