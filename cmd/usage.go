package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var usageCmd = &cobra.Command{
	Use:   "usage <name>",
	Short: "Get usage information for a MCP tool",
	Args:  cobra.ExactArgs(1),
	RunE:  runGetToolUsage,
}

func init() {
	rootCmd.AddCommand(usageCmd)
}

func runGetToolUsage(cmd *cobra.Command, args []string) error {
	t, err := apiClient.GetTool(args[0])
	if err != nil {
		return fmt.Errorf("failed to get tool '%s': %w", args[0], err)
	}
	fmt.Println(t.InputSchema)
	return nil
}
