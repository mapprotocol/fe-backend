package resp

const (
	CodeSuccess             = 2000
	CodeParameterErr        = 4000
	CodeInternalServerError = 5000
)

const (
	CodeExternalServerError = iota + 50001
	CodeButterServerError
	CodeButterNotAvailableRoute
	CodeTONRouteServerError
	CodeTONRouteNotAvailableRoute
)

const (
	CodeOrderNotFound = iota + 4001
	CodeAmountTooFew
	CodeInsufficientLiquidity
)

const (
	MsgSuccess             = "Success"
	MsgInternalServerError = "Internal Server Error"
	MsgParameterErr        = "Invalid Parameter"
)

const (
	MsgButterServerError         = "Butter Server Error"
	MsgButterNotAvailableRoute   = "Butter Not Available Route"
	MsgTONRouteServerError       = "Ton Router Server Error"
	MsgTONRouteNotAvailableRoute = "Ton Router Not Available Route"
)

const (
	MsgOrderNotFound         = "order not found"
	MsgAmountTooFew          = "exchange amount too few"
	MsgInsufficientLiquidity = "insufficient liquidity"
)

var code2msg = map[int]string{
	CodeSuccess:                   MsgSuccess,
	CodeParameterErr:              MsgParameterErr,
	CodeInternalServerError:       MsgInternalServerError,
	CodeButterServerError:         MsgButterServerError,
	CodeButterNotAvailableRoute:   MsgButterNotAvailableRoute,
	CodeTONRouteServerError:       MsgTONRouteServerError,
	CodeTONRouteNotAvailableRoute: MsgTONRouteNotAvailableRoute,
	CodeOrderNotFound:             MsgOrderNotFound,
	CodeAmountTooFew:              MsgAmountTooFew,
	CodeInsufficientLiquidity:     MsgInsufficientLiquidity,
}
