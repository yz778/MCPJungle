package cmd

import (
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
	// TODO: Move the logic of adding query params inside constructURL()
	url := constructURL("/tools")
	if listToolsCmdServerName != "" {
		url += "?server=" + listToolsCmdServerName
	}

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	io.Copy(os.Stdout, resp.Body)
	return nil
}

func runListServers(cmd *cobra.Command, args []string) error {
	resp, err := http.Get(constructURL("/servers"))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	io.Copy(os.Stdout, resp.Body)
	return nil
}
