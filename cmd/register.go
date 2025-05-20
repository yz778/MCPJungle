package cmd

import (
	"bytes"
	"encoding/json"
	"github.com/spf13/cobra"
	"io"
	"net/http"
	"os"
)

var (
	registerCmdServerName string
	registerCmdServerURL  string
	registerCmdServerDesc string
)

var registerMCPServerCmd = &cobra.Command{
	Use:   "register",
	Short: "Register an MCP Server",
	Long:  "Register a MCP Server with the registry.\nA server name is unique across the registry and must not contain a slash '/'",
	RunE:  runRegisterMCPServer,
}

func init() {
	registerMCPServerCmd.Flags().StringVar(
		&registerCmdServerName,
		"name",
		"",
		"MCP server name. If not supplied, name is read from Server metadata",
	)
	registerMCPServerCmd.Flags().StringVar(
		&registerCmdServerURL,
		"url",
		"",
		"URL of the MCP server (eg- http://localhost:8000/mcp)",
	)
	registerMCPServerCmd.Flags().StringVar(
		&registerCmdServerDesc,
		"description",
		"",
		"Server description",
	)

	_ = registerMCPServerCmd.MarkFlagRequired("url")

	rootCmd.AddCommand(registerMCPServerCmd)
}

func runRegisterMCPServer(cmd *cobra.Command, args []string) error {
	payload := map[string]any{
		"name": registerCmdServerName, "url": registerCmdServerURL, "description": registerCmdServerDesc,
	}
	body, _ := json.Marshal(payload)
	resp, err := http.Post(registryServerURL+"/servers", "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	io.Copy(os.Stdout, resp.Body)
	return nil
}
