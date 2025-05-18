package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/duaraghav8/mcpjungle/internal/server"
	"github.com/spf13/cobra"
)

var (
	baseURL string
	port    string
)

func main() {
	root := &cobra.Command{Use: "mcpjungle", Short: "MCP Tool catalogue"}

	// Global flags
	root.PersistentFlags().StringVar(
		&baseURL,
		"registry",
		"http://127.0.0.1:8080",
		"registry server base URL",
	)

	srv := &cobra.Command{
		Use:   "server",
		Short: "Start the registry server",
		Run: func(cmd *cobra.Command, args []string) {
			server.Start(port)
		},
	}
	srv.Flags().StringVar(&port, "port", "8080", "port to bind (overrides $PORT)")

	srvClient := []*cobra.Command{
		cmdRegister(),
		cmdTools(),
		cmdRemove(),
		cmdInvoke(),
		cmdTest(),
	}

	root.AddCommand(srv)
	for _, c := range srvClient {
		root.AddCommand(c)
	}

	if err := root.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func cmdRegister() *cobra.Command {
	var name, url, desc, ttype string
	var tags []string
	c := &cobra.Command{
		Use:   "register",
		Short: "Register a tool in the registry",
		RunE: func(cmd *cobra.Command, args []string) error {
			payload := map[string]any{"name": name, "url": url, "type": ttype, "description": desc, "tags": tags}
			body, _ := json.Marshal(payload)
			resp, err := http.Post(baseURL+"/tools", "application/json", bytes.NewReader(body))
			if err != nil {
				return err
			}
			defer resp.Body.Close()
			io.Copy(os.Stdout, resp.Body)
			return nil
		},
	}
	c.Flags().StringVar(&name, "name", "", "tool name (unique)")
	c.Flags().StringVar(&url, "url", "", "tool base URL")
	c.Flags().StringVar(&ttype, "type", "mcp_server", "tool type")
	c.Flags().StringVar(&desc, "description", "", "description")
	c.Flags().StringSliceVar(&tags, "tag", nil, "tags (repeatable)")
	_ = c.MarkFlagRequired("name")
	_ = c.MarkFlagRequired("url")
	return c
}

func cmdTools() *cobra.Command {
	return &cobra.Command{
		Use:   "tools",
		Short: "List registered tools",
		RunE: func(cmd *cobra.Command, args []string) error {
			resp, err := http.Get(baseURL + "/tools")
			if err != nil {
				return err
			}
			defer resp.Body.Close()
			io.Copy(os.Stdout, resp.Body)
			return nil
		},
	}
}

func cmdRemove() *cobra.Command {
	return &cobra.Command{
		Use:   "remove <name>",
		Short: "Delete a tool by name",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			req, _ := http.NewRequest(http.MethodDelete, baseURL+"/tools/"+args[0], nil)
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				return err
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusNoContent {
				io.Copy(os.Stdout, resp.Body)
			}
			return nil
		},
	}
}

func cmdInvoke() *cobra.Command {
	var input string
	c := &cobra.Command{
		Use:   "invoke <name>",
		Short: "Invoke a tool with JSON input",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			resp, err := http.Post(baseURL, "application/json", bytes.NewReader([]byte(input)))
			if err != nil {
				return err
			}
			defer resp.Body.Close()
			io.Copy(os.Stdout, resp.Body)
			return nil
		},
	}
	c.Flags().StringVar(&input, "input", "{}", "JSON payload")
	return c
}

func cmdTest() *cobra.Command {
	var url string
	c := &cobra.Command{
		Use:   "test",
		Short: "Fetch /.well-known/mcp and print result",
		RunE: func(cmd *cobra.Command, args []string) error {
			resp, err := http.Get(url + "/.well-known/mcp")
			if err != nil {
				return err
			}
			defer resp.Body.Close()
			if resp.StatusCode != 200 {
				return fmt.Errorf("non‑200 status: %d", resp.StatusCode)
			}
			io.Copy(os.Stdout, resp.Body)
			return nil
		},
	}
	c.Flags().StringVar(&url, "url", "", "tool URL to test")
	_ = c.MarkFlagRequired("url")
	return c
}
