package reqerror

type ErrorOption func(*ExternalRequestError)

func WithPath(path string) ErrorOption {
	return func(ere *ExternalRequestError) {
		ere.Path = path
	}
}

func WithMethod(method string) ErrorOption {
	return func(ere *ExternalRequestError) {
		ere.Method = method
	}
}

func WithParams(param string) ErrorOption {
	return func(ere *ExternalRequestError) {
		ere.Param = param
	}
}

func WithCode(code string) ErrorOption {
	return func(ere *ExternalRequestError) {
		ere.Code = code
	}
}

func WithMessage(message string) ErrorOption {
	return func(ere *ExternalRequestError) {
		ere.Message = message
	}

}

func WithError(err error) ErrorOption {
	return func(ere *ExternalRequestError) {
		ere.Err = err
	}
}

func WithPublicError(err string) ErrorOption {
	return func(ere *ExternalRequestError) {
		ere.PublicErr = err
	}
}
