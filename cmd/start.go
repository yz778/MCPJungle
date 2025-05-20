package cmd

import (
	"github.com/duaraghav8/mcpjungle/internal/server"
	"github.com/spf13/cobra"
)

var startServerCmdBindPort string

var startServerCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the MCP registry server",
	Run:   runStartServer,
}

func init() {
	startServerCmd.Flags().StringVar(
		&startServerCmdBindPort,
		"port",
		"8080",
		"port to bind to (overrides $PORT)",
	)
	rootCmd.AddCommand(startServerCmd)
}

func runStartServer(cmd *cobra.Command, args []string) {
	server.Start(startServerCmdBindPort)
}
