package api

import (
	"bytes"
	"io"
	"net/http"

	"github.com/duaraghav8/mcpjungle/internal/service"
	"github.com/gin-gonic/gin"
)

// Invoke forwards the JSON body to the tool URL and streams response back.
func Invoke(c *gin.Context) {
	name := c.Param("name")
	tool, err := service.GetTool(name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "tool not found"})
		return
	}

	reqBody, _ := io.ReadAll(c.Request.Body)
	resp, err := http.Post(tool.URL, "application/json", bytes.NewReader(reqBody))
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	c.Status(resp.StatusCode)
	io.Copy(c.Writer, resp.Body)
}
