package resp

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var EmptyStruct = struct{}{}

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type ListData struct {
	Page  int         `json:"page"`
	Total int64       `json:"total"`
	Items interface{} `json:"items"`
}

func SuccessNil(c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		Code: CodeSuccess,
		Msg:  MsgSuccess,
		Data: EmptyStruct,
	})
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code: CodeSuccess,
		Msg:  MsgSuccess,
		Data: data,
	})
}

func SuccessList(c *gin.Context, page int, total int64, list interface{}) {
	c.JSON(http.StatusOK, Response{
		Code: CodeSuccess,
		Msg:  MsgSuccess,
		Data: ListData{
			Page:  page,
			Total: total,
			Items: list,
		},
	})
}

func Error(c *gin.Context, code int) {
	c.JSON(http.StatusOK, Response{
		Code: code,
		Msg:  code2msg[code],
		Data: EmptyStruct,
	})
}

func ParameterErr(c *gin.Context, msg string) {
	if utils.IsEmpty(msg) {
		msg = code2msg[CodeParameterErr]
	}
	c.JSON(http.StatusOK, Response{
		Code: CodeParameterErr,
		Msg:  msg,
		Data: EmptyStruct,
	})
}
