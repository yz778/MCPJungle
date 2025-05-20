package cmd

import (
	"bytes"
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
	resp, err := http.Post(
		registryServerURL+"/invoke/"+args[0],
		"application/json",
		bytes.NewReader([]byte(invokeCmdInput)),
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	io.Copy(os.Stdout, resp.Body)
	return nil
}
