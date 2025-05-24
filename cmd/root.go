package cmd

import (
	"errors"
	"fmt"
	"github.com/duaraghav8/mcpjungle/client"
	"github.com/duaraghav8/mcpjungle/internal/server"
	"github.com/spf13/cobra"
	"net/http"
	"net/url"
)

// SilentErr is a sentinel error used to indicate that the command should not print an error message
// This is useful when we handle error printing internally but want main to exit with a non-zero status.
// See https://github.com/spf13/cobra/issues/914#issuecomment-548411337
var SilentErr = errors.New("SilentErr")

var registryServerURL string

// apiClient is the global API client used by command handlers to interact with the MCPJungle registry server.
// It is not the best choice to rely on a global variable, but cobra doesn't seem to provide any neat way to
// pass an object down the command tree.
var apiClient *client.Client

var rootCmd = &cobra.Command{
	Use:   "mcpjungle",
	Short: "MCP tool catalog",

	SilenceErrors: true,
	SilenceUsage:  true,

	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
}

func Execute() error {
	// only print usage and error messages if the command usage is incorrect
	rootCmd.SetFlagErrorFunc(func(cmd *cobra.Command, err error) error {
		cmd.Println(err)
		cmd.Println(cmd.UsageString())
		return SilentErr
	})

	rootCmd.PersistentFlags().StringVar(
		&registryServerURL,
		"registry",
		fmt.Sprintf("http://127.0.0.1:%s", BindPortDefault),
		"Base URL of the MCPJungle registry server",
	)

	// Initialize the API client with the registry server URL
	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		apiClient = client.New(registryServerURL, http.DefaultClient)
	}

	return rootCmd.Execute()
}

// constructAPIEndpoint constructs the full API endpoint URL by joining the registry server URL
// with the given suffix path.
func constructAPIEndpoint(suffixPath string) string {
	u, _ := url.JoinPath(registryServerURL, server.ApiV0PathPrefix, suffixPath)
	return u
}
