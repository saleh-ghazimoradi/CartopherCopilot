package jsonrpc

const (
	ErrorParseError     = -32700 // Invalid JSON
	ErrorInvalidRequest = -32600 // Invalid Request object
	ErrorMethodNotFound = -32601 // Method does not exist
	ErrorInvalidParams  = -32602 // Invalid method parameters
	ErrorInternal       = -32603 // Internal JSON-RPC error
)
