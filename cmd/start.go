package cmd

import (
	"fmt"
	"github.com/duaraghav8/mcpjungle/internal/api"
	"github.com/duaraghav8/mcpjungle/internal/db"
	"github.com/duaraghav8/mcpjungle/internal/migrations"
	"github.com/joho/godotenv"
	"github.com/mark3labs/mcp-go/server"
	"github.com/spf13/cobra"
	"os"
)

const (
	BindPortEnvVar  = "PORT"
	BindPortDefault = "8080"
)

var startServerCmdBindPort string

var startServerCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the MCP registry server",
	RunE:  runStartServer,
}

func init() {
	startServerCmd.Flags().StringVar(
		&startServerCmdBindPort,
		"port",
		"",
		fmt.Sprintf("port to bind the server to (overrides env var %s)", BindPortEnvVar),
	)
	rootCmd.AddCommand(startServerCmd)
}

func runStartServer(cmd *cobra.Command, args []string) error {
	_ = godotenv.Load()

	dsn := os.Getenv("DATABASE_URL")
	db.Init(dsn)

	if err := migrations.Migrate(); err != nil {
		return fmt.Errorf("failed to run migrations: %v", err)
	}

	// determine the port to bind the server to
	port := startServerCmdBindPort
	if port == "" {
		port = os.Getenv(BindPortEnvVar)
	}
	if port == "" {
		port = BindPortDefault
	}

	// create the MCP proxy server
	mcpProxyServer := server.NewMCPServer(
		"MCPJungle Proxy MCP Server",
		"0.0.1",
		server.WithToolCapabilities(true),
	)

	// create the API server
	s, err := api.NewServer(port, mcpProxyServer)
	if err != nil {
		return fmt.Errorf("failed to create server: %w", err)
	}

	fmt.Printf("MCPJungle server listening on :%s", port)
	if err := s.Start(); err != nil {
		return fmt.Errorf("failed to run the server: %w", err)
	}

	return nil
}
