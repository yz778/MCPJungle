package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var initServerCmd = &cobra.Command{
	Use:   "init-server",
	Short: "Initialize the MCPJungle Server (for Production Mode only)",
	Long: "If the MCPJungle Server was started in Production Mode, use this command to initialize the server.\n" +
		"Initialization is required before you can use the server.\n",
	RunE: runInitServer,
}

func init() {
	rootCmd.AddCommand(initServerCmd)
}

func runInitServer(cmd *cobra.Command, args []string) error {
	fmt.Println("Initializing the MCPJungle Server in Production Mode...")
	err := apiClient.InitServer()
	if err != nil {
		return fmt.Errorf("failed to initialize the server: %w", err)
	}
	fmt.Println("Done!")
	return nil
}
