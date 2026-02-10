package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/saleh-ghazimoradi/CartopherCopilot/internal/jsonrpc"
	"log/slog"
)

const (
	ProtocolVersion = "2025-11-25"
	ServerName      = "CartopherCopilot"
	ServerVersion   = "1.0.0"
)

type Server struct {
	rpcServer    *jsonrpc.Server
	toolRegistry *Registry
	logger       *slog.Logger
}

func (s *Server) registerHandlers() {
	s.rpcServer.RegisterMethod("initialize", s.handleInitialize)
	s.rpcServer.RegisterMethod("initialized", s.handleInitialized)
	s.rpcServer.RegisterMethod("tools/list", s.handleToolsList)
	s.rpcServer.RegisterMethod("tools/call", s.handleToolsCall)
}

func (s *Server) handleInitialize(params json.RawMessage) (any, error) {
	var req InitializeRequest
	if err := json.Unmarshal(params, &req); err != nil {
		return nil, jsonrpc.NewInvalidParamsError("Invalid initialize parameters")
	}

	s.logger.Info("Client initialized", "client", req.ClientInfo.Name, "clientVersion", req.ClientInfo.Version, "protocolVersion", req.ProtocolVersion)

	return InitializeResult{
		ProtocolVersion: ProtocolVersion,
		Capabilities: ServerCapabilities{
			Tools: &ToolsCapability{
				ListChanged: false,
			},
		},
		ServerInfo: ServerInfo{
			Name:    ServerName,
			Version: ServerVersion,
		},
	}, nil
}

func (s *Server) handleInitialized(params json.RawMessage) (any, error) {
	s.logger.Info("Initialization completed")
	return nil, nil
}

func (s *Server) handleToolsList(params json.RawMessage) (any, error) {
	tools := s.toolRegistry.ListTools()

	s.logger.Debug("Listing tools", "count", len(tools))

	return ToolsListResult{
		Tools: tools,
	}, nil
}

func (s *Server) handleToolsCall(params json.RawMessage) (interface{}, error) {
	var req CallToolRequest

	if err := json.Unmarshal(params, &req); err != nil {
		return nil, jsonrpc.NewInvalidParamsError("Invalid tool call parameters")
	}

	s.logger.Info("Calling tool", "tool", req.Name, "args", req.Arguments)

	ctx := context.Background()

	result, err := s.toolRegistry.ExecuteTool(ctx, req.Name, req.Arguments)
	if err != nil {
		s.logger.Error("Tool execution failed", "err", err)

		return CallToolResult{
			Content: []Content{
				{
					Type: "text",
					Text: fmt.Sprintf("Error: %s", err.Error()),
				},
			},
			IsError: true,
		}, nil
	}

	return result, nil
}

func (s *Server) Start() error {
	return s.rpcServer.ServeStdio()
}

func NewServer(toolRegistry *Registry, logger *slog.Logger) *Server {
	rpcServer := jsonrpc.NewServer(logger)
	server := &Server{
		rpcServer:    rpcServer,
		toolRegistry: toolRegistry,
		logger:       logger,
	}
	server.registerHandlers()
	return server
}
