package resp

const (
	CodeSuccess             = 2000
	CodeParameterErr        = 4000
	CodeInternalServerError = 5000
)

const (
	MsgSuccess             = "Success"
	MsgInternalServerError = "Internal Server Error"
	MsgParameterErr        = "Invalid Parameter"
)

var code2msg = map[int]string{
	CodeSuccess:             MsgSuccess,
	CodeParameterErr:        MsgParameterErr,
	CodeInternalServerError: MsgInternalServerError,
}
