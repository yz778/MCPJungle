package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"net/http"
	"os"
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
	invokeToolCmd.Flags().StringVar(&invokeCmdInput, "input", "{}", "JSON payload")
	rootCmd.AddCommand(invokeToolCmd)
}

func runInvokeTool(cmd *cobra.Command, args []string) error {
	// add tool name to the payload
	var payload map[string]any
	if err := json.Unmarshal([]byte(invokeCmdInput), &payload); err != nil {
		return fmt.Errorf("invalid JSON payload: %w", err)
	}
	payload["name"] = args[0]

	body, _ := json.Marshal(payload)
	u := constructURL("/tools/invoke")
	resp, err := http.Post(u, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("request to server failed: %w", err)
	}
	defer resp.Body.Close()

	io.Copy(os.Stdout, resp.Body)
	return nil
}
