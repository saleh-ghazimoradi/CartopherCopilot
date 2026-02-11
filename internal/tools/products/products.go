package products

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/saleh-ghazimoradi/CartopherCopilot/internal/client"
	"github.com/saleh-ghazimoradi/CartopherCopilot/internal/mcp"
	"log/slog"
)

type ProductToolset struct {
	reg        *mcp.Registry
	logger     *slog.Logger
	restClient *client.RestClient
}

func (r *ProductToolset) registerProductTools() {
	r.reg.Register(mcp.Tool{
		Name:        "list_products",
		Description: "List all available products from the store",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.Property{
				"limit": {
					Type:        "number",
					Description: "Maximum number of products to return (default: 20)",
				},
				"offset": {
					Type:        "number",
					Description: "Number of products to skip (default: 0)",
				},
			},
			Required: []string{},
		},
	}, r.handleListProducts)

	r.reg.Register(mcp.Tool{
		Name:        "search_products",
		Description: "Search for products using a query string",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.Property{
				"q": {
					Type:        "string",
					Description: "Search query to filter products",
				},
				"limit": {
					Type:        "number",
					Description: "Maximum number of products to return (default: 20)",
				},
				"offset": {
					Type:        "number",
					Description: "Number of products to skip (default: 0)",
				},
				"min_price": {
					Type:        "number",
					Description: "Minimum price to filter products",
				},
				"max_price": {
					Type:        "number",
					Description: "Maximum price to filter products",
				},
				"category_id": {
					Type:        "string",
					Description: "The category ID to filter products",
				},
			},
			Required: []string{},
		},
	}, r.searchProducts)

	r.reg.Register(mcp.Tool{
		Name:        "get_product_details",
		Description: "Get detailed information about a specific product by its ID",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.Property{
				"product_id": {
					Type:        "string",
					Description: "The unique identifier of the product",
				},
			},
			Required: []string{"product_id"},
		},
	}, r.getProductDetails)

}

func (r *ProductToolset) getProductDetails(_ context.Context, args map[string]any) (mcp.CallToolResult, error) {
	productID, ok := args["product_id"].(string)
	if !ok || productID == "" {
		return mcp.CallToolResult{}, errors.New("product_id is required and must be a string of numbers. e.g 123")
	}

	response, err := r.restClient.Get(fmt.Sprintf("/products/%s", productID), nil)
	if err != nil {
		return mcp.CallToolResult{}, fmt.Errorf("failed to fetch product details: %w", err)
	}

	var product ProductDetailResponse
	if err := json.Unmarshal(response, &product); err != nil {
		return mcp.CallToolResult{}, fmt.Errorf("failed to parse product detail data: %w", err)
	}

	return mcp.CallToolResult{
		Content: []mcp.Content{
			{
				Type: "text",
				Text: formatProductDetail(product.Data),
			},
		},
	}, nil
}

func (r *ProductToolset) searchProducts(_ context.Context, args map[string]any) (mcp.CallToolResult, error) {

	q, _ := args["q"].(string)
	minPrice, _ := args["min_price"].(float64)
	maxPrice, _ := args["max_price"].(float64)
	categoryId, _ := args["category_id"].(string)

	limit := 20
	offset := 0

	if l, ok := args["limit"].(float64); ok {
		limit = int(l)
	}

	if o, ok := args["offset"].(float64); ok {
		offset = int(o)
	}

	params := map[string]string{
		"limit":  fmt.Sprintf("%d", limit),
		"offset": fmt.Sprintf("%d", offset),
		"q":      q,
	}

	if minPrice > 0 {
		params["min_price"] = fmt.Sprintf("%.2f", minPrice)
	}

	if maxPrice > 0 {
		params["max_price"] = fmt.Sprintf("%.2f", maxPrice)
	}

	if categoryId != "" {
		params["category"] = categoryId
	}

	response, err := r.restClient.Get("/search", params)
	if err != nil {
		return mcp.CallToolResult{}, fmt.Errorf("failed to search products: %w", err)
	}

	r.logger.Info("Fetched products", "response", string(response))
	var products ProductResponse
	if err := json.Unmarshal(response, &products); err != nil {
		return mcp.CallToolResult{}, fmt.Errorf("failed to parse products: %w", err)
	}

	resultText := fmt.Sprintf("Found %d products:\n\n", len(products.Data))
	for i, product := range products.Data {
		resultText += fmt.Sprintf("%d. %s\n", i+1, formatProduct(product))
	}

	return mcp.CallToolResult{
		Content: []mcp.Content{
			{
				Type: "text",
				Text: resultText,
			},
		},
	}, nil

}

func (r *ProductToolset) handleListProducts(_ context.Context, args map[string]any) (mcp.CallToolResult, error) {

	limit := 20
	offset := 0

	if l, ok := args["limit"].(float64); ok {
		limit = int(l)
	}

	if o, ok := args["offset"].(float64); ok {
		offset = int(o)
	}

	params := map[string]string{
		"limit":  fmt.Sprintf("%d", limit),
		"offset": fmt.Sprintf("%d", offset),
	}

	response, err := r.restClient.Get("/products", params)
	if err != nil {
		return mcp.CallToolResult{}, fmt.Errorf("failed to fetch products: %w", err)
	}

	r.logger.Info("Fetched products", "response", string(response))

	var products ProductResponse
	if err := json.Unmarshal(response, &products); err != nil {
		return mcp.CallToolResult{}, fmt.Errorf("failed to parse products: %w", err)
	}

	resultText := fmt.Sprintf("Found %d products:\n\n", len(products.Data))
	for i, product := range products.Data {
		resultText += fmt.Sprintf("%d. %s\n", i+1, formatProduct(product))
	}

	return mcp.CallToolResult{
		Content: []mcp.Content{
			{
				Type: "text",
				Text: resultText,
			},
		},
	}, nil
}

func formatProduct(product Product) string {
	name := product.Name
	price := product.Price
	id := product.Id

	return fmt.Sprintf("**%s** (ID: %d) - $%.2f", name, id, price)
}

func formatProductDetail(product Product) string {
	return fmt.Sprintf(`**Product Details**

ID: %d
Name: %s
Category: %s
Price: $%.2f
Stock: %d units

Description:
%s
`, product.Id, product.Name, product.Category.Name, product.Price, product.Stock, product.Description)
}

func NewProductToolSet(reg *mcp.Registry, restClient *client.RestClient, logger *slog.Logger) *ProductToolset {
	pt := &ProductToolset{
		reg:        reg,
		restClient: restClient,
		logger:     logger,
	}

	pt.registerProductTools()

	return pt
}
