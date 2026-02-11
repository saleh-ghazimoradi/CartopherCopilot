package mcp

func NewToolCallError(message string) CallToolResult {
	return CallToolResult{
		Content: []Content{
			{
				Type: "text",
				Text: message,
			},
		},
		IsError: true,
	}
}
