package cmd

import (
	"github.com/duaraghav8/mcpjungle/internal/server"
	"net/url"
)

func constructAPIEndpoint(suffixPath string) string {
	u, _ := url.JoinPath(registryServerURL, server.ApiV0PathPrefix, suffixPath)
	return u
}

// TODO: Replace all API calls in cmd with calls to an API client SDK
