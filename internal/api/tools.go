package api

import (
	"encoding/json"
	"github.com/duaraghav8/mcpjungle/internal/models"
	"net/http"

	"github.com/duaraghav8/mcpjungle/internal/service"
	"github.com/gin-gonic/gin"
)

func ListToolsHandler(c *gin.Context) {
	server := c.Query("server")
	var (
		tools []models.Tool
		err   error
	)
	if server == "" {
		// no server specified, list all tools
		tools, err = service.ListTools()
	} else {
		// server specified, list tools for that server
		tools, err = service.ListToolsByServer(server)
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tools)
}

// InvokeToolHandler forwards the JSON body to the tool URL and streams response back.
func InvokeToolHandler(c *gin.Context) {
	var args map[string]any
	if err := json.NewDecoder(c.Request.Body).Decode(&args); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "failed to decode request body: " + err.Error()},
		)
		return
	}

	rawName, ok := args["name"]
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing 'name' field in request body"})
		return
	}
	name, ok := rawName.(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "'name' field must be a string"})
		return
	}

	resp, err := service.InvokeTool(c, name, args)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to invoke tool: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}
