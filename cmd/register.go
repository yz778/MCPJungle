package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
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
		"MCP server name",
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

	// TODO: name should not be mandatory.
	//  If not supplied, name should be read from MCP server metadata by the registry.
	_ = registerMCPServerCmd.MarkFlagRequired("name")
	_ = registerMCPServerCmd.MarkFlagRequired("url")

	rootCmd.AddCommand(registerMCPServerCmd)
}

func runRegisterMCPServer(cmd *cobra.Command, args []string) error {
	s, err := apiClient.RegisterServer(registerCmdServerName, registerCmdServerURL, registerCmdServerDesc)
	if err != nil {
		return fmt.Errorf("failed to register server: %w", err)
	}
	fmt.Printf("Server %s registered successfully!\n", s.Name)

	tools, err := apiClient.ListTools(s.Name)
	if err != nil {
		// if we fail to fetch tool list, fail silently because this is not a must-have output
		return nil
	}
	fmt.Println()
	fmt.Println("The following tools are now available from this server:")
	for _, tool := range tools {
		fmt.Printf("- %s: %s\n", tool.Name, tool.Description)
	}

	return nil
}
