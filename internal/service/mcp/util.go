package mcp

import (
	"context"
	"errors"
	"fmt"
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/client/transport"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mcpjungle/mcpjungle/internal/model"
	"net"
	"net/url"
	"regexp"
	"strings"
	"syscall"
)

const serverToolNameSep = "__"

// Only allow letters, numbers, hyphens, and underscores
var validServerName = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

// validServerName checks if the server name is valid.
// Server name must not contain slashes '/'
// Tools in mcpjungle are identified by `<server_name>/<tool_name>` (eg- `github/git_commit`)
// When a tool is invoked, the text before the first slash is treated as the server name.
// eg- In `aws/ec2/create_sg`, `aws` is the MCP server's name and `ec2/create_sg` is the tool.
func validateServerName(name string) error {
	if !validServerName.MatchString(name) {
		return fmt.Errorf("invalid server name: '%s' must not contain slashes or special characters", name)
	}
	return nil
}

// mergeServerToolNames combines the server name and tool name into a single tool name unique across the registry.
func mergeServerToolNames(s, t string) string {
	return s + serverToolNameSep + t
}

// splitServerToolName splits the unique tool name into server name and tool name.
func splitServerToolName(name string) (string, string, bool) {
	serverName, toolName, ok := strings.Cut(name, serverToolNameSep)
	if !ok {
		// there is no separator in tool name, we cannot extract mcp server name
		// this is invalid input
		return "", "", false
	}
	return serverName, toolName, true
}

// isLoopbackURL returns true if rawURL resolves to a loopback address.
// It assumes that rawURL is a valid URL.
func isLoopbackURL(rawURL string) bool {
	u, err := url.Parse(rawURL)
	if err != nil {
		return false // invalid URL, cannot determine loopback
	}
	host := u.Hostname()

	if host == "" {
		return false // no host, not a loopback
	}
	if strings.EqualFold(host, "localhost") {
		return true
	}
	if ip := net.ParseIP(host); ip != nil {
		return ip.IsLoopback()
	}

	return false
}

// createMcpServerConn creates a new MCP server connection and returns the client.
func createMcpServerConn(ctx context.Context, s *model.McpServer) (*client.Client, error) {
	var opts []transport.StreamableHTTPCOption
	if s.BearerToken != "" {
		// If bearer token is provided, set the Authorization header
		o := transport.WithHTTPHeaders(map[string]string{
			"Authorization": "Bearer " + s.BearerToken,
		})
		opts = append(opts, o)
	}

	c, err := client.NewStreamableHttpClient(s.URL, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create streamable HTTP client for MCP server: %w", err)
	}

	initRequest := mcp.InitializeRequest{}
	initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initRequest.Params.ClientInfo = mcp.Implementation{
		Name:    "mcpjungle mcp client for " + s.URL,
		Version: "0.1",
	}
	initRequest.Params.Capabilities = mcp.ClientCapabilities{}

	_, err = c.Initialize(ctx, initRequest)
	if err != nil {
		if errors.Is(err, syscall.ECONNREFUSED) && isLoopbackURL(s.URL) {
			return nil, fmt.Errorf(
				"connection to the MCP server %s was refused. "+
					"If mcpjungle is running inside Docker, use 'host.docker.internal' as your MCP server's hostname",
				s.URL,
			)
		}
		return nil, fmt.Errorf("failed to initialize connection with MCP server: %w", err)
	}

	return c, nil
}
