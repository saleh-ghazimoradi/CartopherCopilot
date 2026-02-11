package orders

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/saleh-ghazimoradi/CartopherCopilot/internal/client"
	"github.com/saleh-ghazimoradi/CartopherCopilot/internal/mcp"
	"log/slog"
)

type OrderToolset struct {
	reg        *mcp.Registry
	logger     *slog.Logger
	restClient *client.RestClient
}

func (o *OrderToolset) registerOrderTools() {
	o.reg.Register(mcp.Tool{
		Name:        "place_order",
		Description: "Place a new order with the items in the shopping cart (requires authentication)",
		InputSchema: mcp.InputSchema{
			Type:       "object",
			Properties: map[string]mcp.Property{},
			Required:   []string{},
		},
	}, o.handlePlaceOrder)
}

func (o *OrderToolset) handlePlaceOrder(_ context.Context, _ map[string]any) (mcp.CallToolResult, error) {

	response, err := o.restClient.WithToken().Post("/orders", nil)
	if err != nil {
		o.logger.Error("Failed to place order", "error", err)
		return mcp.NewToolCallError("Failed to place order"), nil
	}

	var orderRes OrderResponse
	if err := json.Unmarshal(response, &orderRes); err != nil {
		o.logger.Error("Failed to unmarshal order response", "error", err)
		return mcp.NewToolCallError("Failed to parse order response"), nil
	}

	return mcp.CallToolResult{
		Content: []mcp.Content{
			{
				Type: "text",
				Text: fmt.Sprintf("Order placed successfully! Order ID: %d, Total Amount: $%.2f",
					orderRes.Data.Id,
					orderRes.Data.Total),
			},
		},
		IsError: false,
	}, nil
}

func NewOrderToolset(reg *mcp.Registry, restClient *client.RestClient, logger *slog.Logger) *OrderToolset {
	ot := &OrderToolset{
		reg:        reg,
		restClient: restClient,
		logger:     logger,
	}

	ot.registerOrderTools()

	return ot
}
