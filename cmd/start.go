package cmd

import (
	"fmt"
	"github.com/duaraghav8/mcpjungle/internal/api"
	"github.com/duaraghav8/mcpjungle/internal/db"
	"github.com/duaraghav8/mcpjungle/internal/migrations"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"log"
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

	s := api.NewServer()

	port := startServerCmdBindPort
	if port == "" {
		port = os.Getenv(BindPortEnvVar)
	}
	if port == "" {
		port = BindPortDefault
	}

	log.Printf("MCPJungle server listening on :%s", port)
	if err := s.Run(":" + port); err != nil {
		return fmt.Errorf("failed to run the server: %v", err)
	}

	return nil
}
