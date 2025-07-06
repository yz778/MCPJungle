package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Server represents an MCP server registered in the MCPJungle registry.
type Server struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	URL         string `json:"url"`
}

// RegisterServerInput is the input structure for registering a new MCP server.
type RegisterServerInput struct {
	Name        string `json:"name"`
	Description string `json:"description"`

	// URL is mandatory and must be a valid http/https URL (eg- https://example.com/mcp).
	// MCPJungle only supports streamable HTTP transport as of now.
	URL string `json:"url"`

	// BearerToken is an optional token used for authenticating requests to the MCP server.
	// It is useful when the upstream MCP server requires static tokens (e.g., API tokens) for authentication.
	BearerToken string `json:"bearer_token,omitempty"`
}

// RegisterServer registers a new MCP server with the registry.
func (c *Client) RegisterServer(server *RegisterServerInput) (*Server, error) {
	u, _ := c.constructAPIEndpoint("/servers")
	body, err := json.Marshal(server)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize server data into JSON: %w", err)
	}

	req, err := c.newRequest(http.MethodPost, u, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request to %s: %w", u, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("request failed with status: %d, message: %s", resp.StatusCode, body)
	}

	var registeredServer Server
	if err := json.NewDecoder(resp.Body).Decode(&registeredServer); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return &registeredServer, nil
}

// ListServers fetches the list of registered servers.
func (c *Client) ListServers() ([]*Server, error) {
	u, _ := c.constructAPIEndpoint("/servers")
	req, err := c.newRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request to %s: %w", u, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("request failed with status: %d, message: %s", resp.StatusCode, body)
	}

	var servers []*Server
	if err := json.NewDecoder(resp.Body).Decode(&servers); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return servers, nil
}

// DeregisterServer deletes a server by name.
func (c *Client) DeregisterServer(name string) error {
	u, _ := c.constructAPIEndpoint("/servers/" + name)
	req, _ := c.newRequest(http.MethodDelete, u, nil)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request to %s: %w", u, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status from server: %s, body: %s", resp.Status, body)
	}
	return nil
}
