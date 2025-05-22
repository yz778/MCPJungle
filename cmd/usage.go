package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"net/http"
	"os"
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
	req, err := http.NewRequest(http.MethodGet, constructAPIEndpoint("/tool"), nil)
	if err != nil {
		return fmt.Errorf("failed to create HTTP request for server: %w", err)
	}

	q := req.URL.Query()
	q.Add("name", args[0])
	req.URL.RawQuery = q.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to fetch information about tool '%s': %w", args[0], err)
	}
	defer resp.Body.Close()

	io.Copy(os.Stdout, resp.Body)
	return nil
}
