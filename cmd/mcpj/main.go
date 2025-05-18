package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var baseURL string

func main() {
	root := &cobra.Command{Use: "mcpj"}
	root.PersistentFlags().StringVar(&baseURL, "registry", "http://localhost:8080", "registry base URL")

	root.AddCommand(cmdServe())
	root.AddCommand(cmdRegister())
	root.AddCommand(cmdTools())
	root.AddCommand(cmdRemove())
	root.AddCommand(cmdInvoke())
	root.AddCommand(cmdTest())

	if err := root.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func cmdServe() *cobra.Command {
	return &cobra.Command{
		Use:   "serve",
		Short: "Start the registry server (wrapper around the same binary)",
		RunE: func(cmd *cobra.Command, args []string) error {
			// re‑exec server main
			return exec.Command(os.Args[0], "server").Run()
		},
	}
}

func cmdRegister() *cobra.Command {
	var name, url, desc, ttype string
	var tags []string
	c := &cobra.Command{
		Use:   "register",
		Short: "Register a tool",
		RunE: func(cmd *cobra.Command, args []string) error {
			payload := map[string]any{
				"name":        name,
				"url":         url,
				"type":        ttype,
				"description": desc,
				"tags":        tags,
			}
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
		Short: "List tools",
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
		Short: "Delete a tool",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			req, _ := http.NewRequest(http.MethodDelete, baseURL+"/tools/"+args[0], nil)
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				return err
			}
			defer resp.Body.Close()
			if resp.StatusCode != 204 {
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
			resp, err := http.Post(baseURL+"/invoke/"+args[0], "application/json", bytes.NewReader([]byte(input)))
			if err != nil {
				return err
			}
			defer resp.Body.Close()
			io.Copy(os.Stdout, resp.Body)
			return nil
		},
	}
	c.Flags().StringVar(&input, "input", "{}", "JSON payload to send")
	return c
}

func cmdTest() *cobra.Command {
	var testURL string
	c := &cobra.Command{
		Use:   "test",
		Short: "Test MCP compliance by fetching /.well-known/mcp",
		RunE: func(cmd *cobra.Command, args []string) error {
			resp, err := http.Get(testURL + "/.well-known/mcp")
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
	c.Flags().StringVar(&testURL, "url", "", "tool URL to test")
	_ = c.MarkFlagRequired("url")
	return c
}
