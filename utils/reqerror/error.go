package reqerror

import (
	"fmt"
	"github.com/mapprotocol/fe-backend/utils"
	"strings"
)

type ExternalRequestError struct {
	Path      string
	Method    string
	Param     string
	Code      string
	Message   string
	Err       error
	PublicErr string
}

func NewExternalRequestError(path string, opts ...ErrorOption) *ExternalRequestError {
	ext := ExternalRequestError{
		Path: path,
	}
	for _, o := range opts {
		o(&ext)
	}
	return &ext
}

func (e *ExternalRequestError) HasPublicError() bool {
	return !utils.IsEmpty(e.PublicErr)
}
func (e *ExternalRequestError) PublicError() string {
	return e.PublicErr
}

func (e *ExternalRequestError) Error() string {
	msg := fmt.Sprintf("[Error while making external request to %s, ", e.Path)

	if e.Method != "" {
		msg += fmt.Sprintf(" method: %s, ", e.Method)
	}
	if e.Param != "" {
		msg += fmt.Sprintf(" param: %s, ", e.Param)
	}
	if e.Code != "" {
		msg += fmt.Sprintf("code: %s, ", e.Code)
	}
	if e.Message != "" {
		msg += fmt.Sprintf("message: %s, ", e.Message)
	}
	if e.Err != nil {
		msg += fmt.Sprintf("error: %s", e.Err.Error())
	}
	msg = strings.TrimSuffix(msg, ", ")
	msg += "]"
	return msg
}
