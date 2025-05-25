package client

import (
	"github.com/duaraghav8/mcpjungle/internal/api"
	"net/http"
	"net/url"
)

// Client represents a client for interacting with the MCPJungle HTTP API
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewClient(baseURL string, httpClient *http.Client) *Client {
	return &Client{
		BaseURL:    baseURL,
		HTTPClient: httpClient,
	}
}

// constructAPIEndpoint constructs the full API endpoint URL where a request must be sent
func (c *Client) constructAPIEndpoint(suffixPath string) (string, error) {
	return url.JoinPath(c.BaseURL, api.V0PathPrefix, suffixPath)
}
