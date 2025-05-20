package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var registryServerURL string

var rootCmd = &cobra.Command{
	Use:   "mcpjungle",
	Short: "MCP tool catalog",
}

func Execute() {
	rootCmd.PersistentFlags().StringVar(
		&registryServerURL,
		"registry",
		"http://127.0.0.1:8080",
		"Base URL of the mcpjungle registry server",
	)

	rootCmd.CompletionOptions.DisableDefaultCmd = true

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
