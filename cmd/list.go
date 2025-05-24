package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"net/http"
	"os"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List resources",
}

var listToolsCmdServerName string

var listToolsCmd = &cobra.Command{
	Use:   "tools",
	Short: "List available tools",
	Long:  "List tools available either from a specific MCP server or across all MCP servers registered in the registry.",
	RunE:  runListTools,
}

var listServersCmd = &cobra.Command{
	Use:   "servers",
	Short: "List registered MCP servers",
	RunE:  runListServers,
}

func init() {
	listToolsCmd.Flags().StringVar(
		&listToolsCmdServerName,
		"server",
		"",
		"Filter tools by server name",
	)

	listCmd.AddCommand(listToolsCmd)
	listCmd.AddCommand(listServersCmd)
	rootCmd.AddCommand(listCmd)
}

func runListTools(cmd *cobra.Command, args []string) error {
	req, _ := http.NewRequest(http.MethodGet, constructAPIEndpoint("/tools"), nil)

	// if server name is provided, add it to the query parameters
	if listToolsCmdServerName != "" {
		q := req.URL.Query()
		q.Add("server", listToolsCmdServerName)
		req.URL.RawQuery = q.Encode()
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to list tools: %w", err)
	}
	defer resp.Body.Close()

	io.Copy(os.Stdout, resp.Body)
	return nil
}

func runListServers(cmd *cobra.Command, args []string) error {
	resp, err := apiClient.ListServers()
	if err != nil {
		return err
	}

	fmt.Println(resp)
	return nil
}
