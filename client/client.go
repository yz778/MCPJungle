package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/mcpjungle/mcpjungle/internal/api"
	"github.com/mcpjungle/mcpjungle/internal/model"
	"io"
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

// InitServer sends a request to initialize the server in production mode
func (c *Client) InitServer() error {
	u, _ := url.JoinPath(c.BaseURL, "/init")

	payload := struct {
		Mode string `json:"mode"`
	}{
		Mode: string(model.ModeProd),
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := c.HTTPClient.Post(u, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to send request to %s: %w", u, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("request failed with status: %d, message: %s", resp.StatusCode, body)
	}

	return nil
}

// constructAPIEndpoint constructs the full API endpoint URL where a request must be sent
func (c *Client) constructAPIEndpoint(suffixPath string) (string, error) {
	return url.JoinPath(c.BaseURL, api.V0PathPrefix, suffixPath)
}
