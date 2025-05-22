package cmd

import (
	"github.com/spf13/cobra"
	"io"
	"net/http"
	"os"
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
		io.Copy(os.Stdout, resp.Body)
	}
	return nil
}
