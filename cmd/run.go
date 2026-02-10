package cmd

import (
	"fmt"
	"github.com/saleh-ghazimoradi/CartopherCopilot/config"
	"github.com/saleh-ghazimoradi/CartopherCopilot/internal/mcp"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("run called")

		logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))

		cfg, err := config.GetConfig()
		if err != nil {
			logger.Error("failed to get config", "error", err.Error())
			os.Exit(1)
		}

		logger.Info("Starting CartopherCopilot Server", "api_url", cfg.APIURL, "auth_token_configured", cfg.AuthToken)

		//restClient := client.NewRestClient(cfg.APIURL, cfg.AuthToken, logger)

		toolRegistry := mcp.NewRegistry(logger)

		logger.Info("Registry tools", "tool_count", len(toolRegistry.ListTools()))

		mcpServer := mcp.NewServer(toolRegistry, logger)

		if err := mcpServer.Start(); err != nil {
			logger.Error("Server error", "error", err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
