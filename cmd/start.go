package cmd

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/mark3labs/mcp-go/server"
	"github.com/mcpjungle/mcpjungle/internal/api"
	"github.com/mcpjungle/mcpjungle/internal/db"
	"github.com/mcpjungle/mcpjungle/internal/migrations"
	"github.com/mcpjungle/mcpjungle/internal/model"
	"github.com/mcpjungle/mcpjungle/internal/service"
	"github.com/mcpjungle/mcpjungle/internal/service/config"
	"github.com/spf13/cobra"
	"os"
)

const (
	BindPortEnvVar  = "PORT"
	BindPortDefault = "8080"
)

var (
	startServerCmdBindPort    string
	startServerCmdProdEnabled bool
)

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
	startServerCmd.Flags().BoolVar(
		&startServerCmdProdEnabled,
		"prod",
		false,
		fmt.Sprintf("Run server in production mode (for enterprises)"),
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
	// Migrations should ideally be decoupled from both the server and the startup phase
	// (should be run as a separate command).
	// However, for the user's convenience, we run them as part of startup command for now.
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

	configService := config.NewServerConfigService(dbConn)

	// create the API server
	s, err := api.NewServer(port, mcpProxyServer, mcpService, configService)
	if err != nil {
		return fmt.Errorf("failed to create server: %v", err)
	}

	// Silently initialize the server if the mode is dev.
	// Individual (dev mode) users need not worry about server initialization.
	if !startServerCmdProdEnabled {
		if err := s.Init(model.ModeDev); err != nil {
			return fmt.Errorf("failed to initialize server in development mode: %v", err)
		}
	} else {
		fmt.Println(
			"Starting server in production mode, don't forget to initialize it by running `init-server`",
		)
	}

	fmt.Printf("MCPJungle server listening on :%s", port)
	if err := s.Start(); err != nil {
		return fmt.Errorf("failed to run the server: %v", err)
	}

	return nil
}
