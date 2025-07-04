package cmd

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/mark3labs/mcp-go/server"
	"github.com/mcpjungle/mcpjungle/internal/api"
	"github.com/mcpjungle/mcpjungle/internal/db"
	"github.com/mcpjungle/mcpjungle/internal/migrations"
	"github.com/mcpjungle/mcpjungle/internal/model"
	"github.com/mcpjungle/mcpjungle/internal/service/config"
	"github.com/mcpjungle/mcpjungle/internal/service/mcp"
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
	Long: "Starts the MCPJungle HTTP registry server and the MCP Proxy server.\n" +
		"The server is started in Development mode by default, which is ideal for individual users.\n",
	RunE: runStartServer,
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
		fmt.Sprintf("Run the server in Production mode (suitable for enterprises)"),
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

	mcpService, err := mcp.NewMCPService(dbConn, mcpProxyServer)
	if err != nil {
		return fmt.Errorf("failed to create MCP service: %v", err)
	}

	configService := config.NewServerConfigService(dbConn)

	// create the API server
	opts := &api.ServerOptions{
		Port:           port,
		MCPProxyServer: mcpProxyServer,
		MCPService:     mcpService,
		ConfigService:  configService,
	}
	s, err := api.NewServer(opts)
	if err != nil {
		return fmt.Errorf("failed to create server: %v", err)
	}

	desiredMode := model.ModeDev
	if startServerCmdProdEnabled {
		desiredMode = model.ModeProd
	}

	ok, err := s.IsInitialized()
	if err != nil {
		return fmt.Errorf("failed to check if server is initialized: %v", err)
	}
	if ok {
		// If the server is already initialized, then the mode supplied to this command (desired mode)
		// must match the configured mode.
		mode, err := s.GetMode()
		if err != nil {
			return fmt.Errorf("failed to get server mode: %v", err)
		}
		if desiredMode != mode {
			return fmt.Errorf(
				"server is already initialized in %s mode, cannot start in %s mode",
				mode, desiredMode,
			)
		}
	} else {
		// If server isn't already initialized and the desired mode is dev, silently initialize the server.
		// Individual (dev mode) users need not worry about server initialization.
		if desiredMode == model.ModeDev {
			if err := s.Init(desiredMode); err != nil {
				return fmt.Errorf("failed to initialize server in development mode: %v", err)
			}
		} else {
			// If desired mode is prod, then server initialization is a manual next step to be taken by the user.
			// This is so that they can obtain the admin access token on their client machine.
			fmt.Println(
				"Starting server in Production mode, don't forget to initialize it by running `init-server`",
			)
		}
	}

	fmt.Printf("MCPJungle HTTP server listening on :%s\n\n", port)
	if err := s.Start(); err != nil {
		return fmt.Errorf("failed to run the server: %v\n", err)
	}

	return nil
}
