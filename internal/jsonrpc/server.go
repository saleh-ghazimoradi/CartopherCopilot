package jsonrpc

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"os"
)

type Handler func(params json.RawMessage) (any, error)

type Server struct {
	handlers map[string]Handler
	logger   *slog.Logger
}

func (s *Server) RegisterMethod(method string, handler Handler) {
	s.handlers[method] = handler
	s.logger.Debug("registered handler", "method", method)
}

func (s *Server) HandleRequest(req *Request) *Response {
	s.logger.Debug("Handling request", "method", req.Method, "id", req.Id)

	if err := req.Validate(); err != nil {
		var jsonErr *Error
		if errors.As(err, &jsonErr) {
			return NewErrorResponse(req.Id, jsonErr)
		}
		return NewErrorResponse(req.Id, NewInternalError(err))
	}

	handler, ok := s.handlers[req.Method]
	if !ok {
		return NewErrorResponse(req.Id, NewMethodNotFoundError(req.Method))
	}

	result, err := handler(req.Params)
	if err != nil {
		var jsonErr *Error
		if errors.As(err, &jsonErr) {
			return NewErrorResponse(req.Id, jsonErr)
		}
		return NewErrorResponse(req.Id, NewInternalError(err))
	}

	return NewSuccessResponse(result, req.Id)
}

func (s *Server) ServeStdio() error {
	s.logger.Info("Starting server on stdio")

	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)

	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				s.logger.Info("EOF received, shutting down server")
				return nil
			}
			s.logger.Error("unable to read from stdin", "error", err)
			return err
		}

		s.logger.Debug("Received request", "request", string(line))

		var req Request

		if err := json.Unmarshal(line, &req); err != nil {
			res := NewErrorResponse(nil, NewParseError(err))
			s.writeResponse(writer, res)
			continue
		}

		res := s.HandleRequest(&req)
		if !req.IsNotification() {
			s.writeResponse(writer, res)
		}
	}
}

func (s *Server) writeResponse(writer *bufio.Writer, res *Response) {
	bs, err := json.Marshal(res)
	if err != nil {
		s.logger.Error("unable to marshal response", "error", err)
		return
	}
	s.logger.Debug("Sending response", "response", string(bs))

	_, _ = writer.Write(bs)
	_ = writer.WriteByte('\n')
	_ = writer.Flush()
}

func NewServer(logger *slog.Logger) *Server {
	return &Server{
		handlers: make(map[string]Handler),
		logger:   logger,
	}
}
