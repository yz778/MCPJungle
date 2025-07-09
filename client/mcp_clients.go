package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// McpClient represents an MCP client that is authorized to access the MCPJungle MCP Proxy server.
type McpClient struct {
	Name        string `json:"name"`
	Description string `json:"description"`

	// AllowList is a comma-separated list of MCP Servers that this client
	// is allowed to access from MCPJungle.
	AllowList []string `json:"allow_list"`
}

func (c *Client) ListMcpClients() ([]McpClient, error) {
	u, _ := c.constructAPIEndpoint("/clients")

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

	var clients []McpClient
	if err := json.NewDecoder(resp.Body).Decode(&clients); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return clients, nil
}

func (c *Client) DeleteMcpClient(name string) error {
	u, _ := c.constructAPIEndpoint("/clients/" + name)

	req, err := c.newRequest(http.MethodDelete, u, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request to %s: %w", u, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("request failed with status: %d, message: %s", resp.StatusCode, body)
	}

	return nil
}

func (c *Client) CreateMcpClient(mcpClient *McpClient) (string, error) {
	u, _ := c.constructAPIEndpoint("/clients")

	body, err := json.Marshal(mcpClient)
	if err != nil {
		return "", fmt.Errorf("failed to marshal client data: %w", err)
	}

	req, err := c.newRequest(http.MethodPost, u, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request to %s: %w", u, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("request failed with status: %d, message: %s", resp.StatusCode, body)
	}

	var response struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	return response.AccessToken, nil
}
