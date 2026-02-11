package main

import (
	"github.com/saleh-ghazimoradi/CartopherCopilot/config"
	"github.com/saleh-ghazimoradi/CartopherCopilot/internal/client"
	"github.com/saleh-ghazimoradi/CartopherCopilot/internal/mcp"
	"github.com/saleh-ghazimoradi/CartopherCopilot/internal/tools/cart"
	"github.com/saleh-ghazimoradi/CartopherCopilot/internal/tools/orders"
	"github.com/saleh-ghazimoradi/CartopherCopilot/internal/tools/products"
	"log/slog"
	"os"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))

	cfg, err := config.GetConfig()
	if err != nil {
		logger.Error("failed to get config", "error", err.Error())
		os.Exit(1)
	}

	logger.Info("Starting CartopherCopilot Server", "api_url", cfg.APIURL, "auth_token_configured", cfg.AuthToken)

	restClient := client.NewRestClient(cfg.APIURL, cfg.AuthToken, logger)

	toolRegistry := mcp.NewRegistry(logger)

	products.NewProductToolSet(toolRegistry, restClient, logger)
	cart.NewCartToolset(toolRegistry, restClient, logger)
	orders.NewOrderToolset(toolRegistry, restClient, logger)

	logger.Info("Registry tools", "tool_count", len(toolRegistry.ListTools()))

	mcpServer := mcp.NewServer(toolRegistry, logger)

	if err := mcpServer.Start(); err != nil {
		logger.Error("Server error", "error", err.Error())
	}
}
