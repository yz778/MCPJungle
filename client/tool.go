package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// ToolInputSchema defines the schema for the input parameters of a tool
type ToolInputSchema struct {
	Type       string         `json:"type"`
	Properties map[string]any `json:"properties,omitempty"`
	Required   []string       `json:"required,omitempty"`
}

// Tool represents a tool provided by an MCP Server registered in the registry.
type Tool struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	InputSchema ToolInputSchema `json:"input_schema"`
}

// ListTools fetches the list of tools, optionally filtered by server name.
func (c *Client) ListTools(server string) ([]*Tool, error) {
	u, _ := c.constructAPIEndpoint("/tools")
	req, _ := http.NewRequest(http.MethodGet, u, nil)
	if server != "" {
		q := req.URL.Query()
		q.Add("server", server)
		req.URL.RawQuery = q.Encode()
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request to %s: %w", req.URL.String(), err)
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("request failed with status: %d, message: %s", resp.StatusCode, body)
	}

	var tools []*Tool
	if err := json.NewDecoder(resp.Body).Decode(&tools); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return tools, nil
}
