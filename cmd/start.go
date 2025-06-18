package cmd

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/mark3labs/mcp-go/server"
	"github.com/mcpjungle/mcpjungle/internal/api"
	"github.com/mcpjungle/mcpjungle/internal/db"
	"github.com/mcpjungle/mcpjungle/internal/migrations"
	"github.com/mcpjungle/mcpjungle/internal/service"
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

	// connect to the DB and run migrations
	dsn := os.Getenv("DATABASE_URL")
	dbConn, err := db.NewDBConnection(dsn)
	if err != nil {
		return err
	}
	if err := migrations.Migrate(dbConn); err != nil {
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

	mcpService, err := service.NewMCPService(dbConn, mcpProxyServer)
	if err != nil {
		return fmt.Errorf("failed to create MCP service: %v", err)
	}

	// create the API server
	s, err := api.NewServer(port, mcpProxyServer, mcpService)
	if err != nil {
		return fmt.Errorf("failed to create server: %v", err)
	}

	fmt.Printf("MCPJungle server listening on :%s", port)
	if err := s.Start(); err != nil {
		return fmt.Errorf("failed to run the server: %v", err)
	}

	return nil
}
