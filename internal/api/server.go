package api

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"log"
)

const V0PathPrefix = "/api/v0"

// NewServer initializes a new Gin server with the MCPJungle MCP server and API endpoints.
func NewServer() *gin.Engine {
	r := gin.Default()

	r.GET(
		"/health",
		func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		},
	)

	// Set up the proxy MCP server on /mcp path
	mcpServer := server.NewMCPServer(
		"MCPJungle Proxy MCP Server",
		"0.0.1",
		server.WithToolCapabilities(true),
	)

	var h = func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		upstreamServer := "http://127.0.0.1:8000/mcp"

		fmt.Println("Received multiply tool call:", request.Method, request.Params)
		a := request.GetInt("a", 0)
		b := request.GetInt("b", 0)
		fmt.Println(request.Params.Name, a, b)

		mcpClient, err := client.NewStreamableHttpClient(upstreamServer)
		if err != nil {
			return nil, fmt.Errorf("failed to create MCP client: %w", err)
		}

		fmt.Println("Initializing client...")
		initRequest := mcp.InitializeRequest{}
		initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
		initRequest.Params.ClientInfo = mcp.Implementation{
			Name:    "MCP client using streamable http",
			Version: "0.0.1",
		}
		initRequest.Params.Capabilities = mcp.ClientCapabilities{}

		_, err = mcpClient.Initialize(context.Background(), initRequest)
		if err != nil {
			log.Fatalf("Failed to initialize client: %v", err)
		}
		fmt.Println("Done initializing client")

		return mcpClient.CallTool(ctx, request)
	}

	t := mcp.NewTool(
		"multiply",
		mcp.WithDescription("Multiplies two numbers"),
		mcp.WithNumber("a", mcp.Description("First number"), mcp.Required()),
		mcp.WithNumber("b", mcp.Description("Second number"), mcp.Required()),
	)
	mcpServer.AddTool(t, h)

	streamableHttpServer := server.NewStreamableHTTPServer(mcpServer)
	r.Any("/mcp", gin.WrapH(streamableHttpServer))

	apiV0 := r.Group(V0PathPrefix)
	{
		apiV0.POST("/servers", registerServerHandler)
		apiV0.DELETE("/servers/:name", deregisterServerHandler)
		apiV0.GET("/servers", listServersHandler)
		apiV0.GET("/tools", listToolsHandler)
		apiV0.POST("/tools/invoke", invokeToolHandler)
		apiV0.GET("/tool", getToolHandler)
	}

	return r
}
