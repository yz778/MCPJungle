package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"net/http"
)

var deregisterMCPServerCmd = &cobra.Command{
	Use:   "deregister",
	Short: "Deregister an MCP Server",
	Long:  "Remove an MCP server from the registry. This also deregisters all tools provided by the server.",
	Args:  cobra.ExactArgs(1),
	RunE:  runDeregisterMCPServer,
}

func init() {
	rootCmd.AddCommand(deregisterMCPServerCmd)
}

func runDeregisterMCPServer(cmd *cobra.Command, args []string) error {
	url := constructAPIEndpoint("/servers/" + args[0])
	req, _ := http.NewRequest(http.MethodDelete, url, nil)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("server responded with unexpected status %s: %s", resp.Status, resp.Body)
	}
	return nil
}
