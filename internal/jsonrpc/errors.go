package jsonrpc

const (
	ErrorParseError     = -32700 // Invalid JSON
	ErrorInvalidRequest = -32600 // Invalid Request object
	ErrorMethodNotFound = -32601 // Method does not exist
	ErrorInvalidParams  = -32602 // Invalid method parameters
	ErrorInternal       = -32603 // Internal JSON-RPC error
)

func NewError(code int, message string, data any) *Error {
	return &Error{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

func NewParseError(data any) *Error {
	return NewError(ErrorParseError, "Parse Error", data)
}

func NewInvalidRequestError(data any) *Error {
	return NewError(ErrorInvalidRequest, "Invalid Request", data)
}

func NewMethodNotFoundError(data any) *Error {
	return NewError(ErrorMethodNotFound, "Method not found", data)
}

func NewInvalidParamsError(data any) *Error {
	return NewError(ErrorInvalidParams, "Invalid params", data)
}

func NewInternalError(data any) *Error {
	return NewError(ErrorInternal, "Internal Error", data)
}
