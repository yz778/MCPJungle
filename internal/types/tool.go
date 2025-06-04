package types

// ToolInvokeResult represents the result of a Tool call.
// It is designed to be passed down to the end user.
type ToolInvokeResult struct {
	// Meta and IsError are taken directly from the MCP-compliant Tool response object.
	Meta    map[string]any `json:"_meta,omitempty"`
	IsError bool           `json:"isError,omitempty"`

	// TextContent contains the text output of the tool invocation.
	// As of now, MCPJungle only supports text output from tools.
	TextContent []string `json:"textContent"`
}
