package cmd

import (
	"fmt"
	"github.com/duaraghav8/mcpjungle/internal/server"
	"net/url"
	"os"

	"github.com/spf13/cobra"
)

var registryServerURL string

var rootCmd = &cobra.Command{
	Use:   "mcpjungle",
	Short: "MCP tool catalog",
}

func Execute() {
	rootCmd.PersistentFlags().StringVar(
		&registryServerURL,
		"registry",
		fmt.Sprintf("http://127.0.0.1:%s", BindPortDefault),
		"Base URL of the mcpjungle registry server",
	)

	rootCmd.CompletionOptions.DisableDefaultCmd = true

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// constructAPIEndpoint constructs the full API endpoint URL by joining the registry server URL
// with the given suffix path.
func constructAPIEndpoint(suffixPath string) string {
	u, _ := url.JoinPath(registryServerURL, server.ApiV0PathPrefix, suffixPath)
	return u
}
