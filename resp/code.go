package resp

const (
	CodeSuccess             = 2000
	CodeParameterErr        = 4000
	CodeInternalServerError = 5000
)

const (
	CodeButterServerError = iota + 50001
	CodeButterNotAvailableRoute
	CodeTONRouteServerError
	CodeTONRouteNotAvailableRoute
)

const (
	CodeOrderNotFound = 4001
)

const (
	MsgSuccess             = "Success"
	MsgInternalServerError = "Internal Server Error"
	MsgParameterErr        = "Invalid Parameter"
)

const (
	MsgOrderNotFound = "order not found"
)

var code2msg = map[int]string{
	CodeSuccess:             MsgSuccess,
	CodeParameterErr:        MsgParameterErr,
	CodeInternalServerError: MsgInternalServerError,
	CodeOrderNotFound:       MsgOrderNotFound,
}
