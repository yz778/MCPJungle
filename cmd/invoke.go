package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
)

var invokeCmdInput string

var invokeToolCmd = &cobra.Command{
	Use:   "invoke <name>",
	Short: "Invoke a tool",
	Long:  "Invokes a tool supplied by a registered MCP server",
	Args:  cobra.ExactArgs(1),
	RunE:  runInvokeTool,
}

func init() {
	invokeToolCmd.Flags().StringVar(&invokeCmdInput, "input", "{}", "valid JSON payload")
	rootCmd.AddCommand(invokeToolCmd)
}

func runInvokeTool(cmd *cobra.Command, args []string) error {
	var input map[string]any
	if err := json.Unmarshal([]byte(invokeCmdInput), &input); err != nil {
		return fmt.Errorf("invalid input: %w", err)
	}

	resp, err := apiClient.InvokeTool(args[0], input)
	if err != nil {
		return fmt.Errorf("failed to invoke tool: %w", err)
	}

	fmt.Println("Response from tool:")
	fmt.Println()
	fmt.Println(resp)
	return nil
}
