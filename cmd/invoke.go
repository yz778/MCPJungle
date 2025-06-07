package cmd

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"time"
)

var invokeCmdInput string

var invokeToolCmd = &cobra.Command{
	Use:   "invoke <name>",
	Short: "Invoke a tool",
	Long:  "Invokes a tool supplied by a registered MCP server",
	Args:  cobra.ExactArgs(1),
	RunE:  runInvokeTool,
}

func init() {
	invokeToolCmd.Flags().StringVar(&invokeCmdInput, "input", "{}", "valid JSON payload")
	rootCmd.AddCommand(invokeToolCmd)
}

func runInvokeTool(cmd *cobra.Command, args []string) error {
	var input map[string]any
	if err := json.Unmarshal([]byte(invokeCmdInput), &input); err != nil {
		return fmt.Errorf("invalid input: %w", err)
	}

	result, err := apiClient.InvokeTool(args[0], input)
	if err != nil {
		return fmt.Errorf("failed to invoke tool: %w", err)
	}

	if result.IsError {
		fmt.Println("The tool returned an error:")
		for k, v := range result.Meta {
			fmt.Printf("%s: %v\n", k, v)
		}
	} else {
		fmt.Println("Response from tool:")
	}

	// result Content needs to be printed regardless of whether the tool returned an error or not
	// because it may contain useful information
	fmt.Println()
	for _, c := range result.Content {
		cType, ok := c["type"]
		if !ok {
			return fmt.Errorf("content item does not have a 'type' field: %v", c)
		}
		switch cType {

		case "text":
			textContent, ok := c["text"].(string)
			if !ok {
				return fmt.Errorf("text content item does not have a 'text' field: %v", c)
			}
			fmt.Println(textContent)

		case "image":
			dataStr, ok := c["data"].(string)
			if !ok {
				return fmt.Errorf("image content item does not have a valid 'data' field: %v", c)
			}
			mimeType, ok := c["mimeType"].(string)
			if !ok {
				return fmt.Errorf("image content item does not have a valid 'mimeType' field: %v", c)
			}

			// Decode base64 image data
			imgData, err := base64.StdEncoding.DecodeString(dataStr)
			if err != nil {
				return fmt.Errorf("failed to decode base64 image data: %w", err)
			}

			// Determine file extension from MIME type
			ext := ".img"
			switch mimeType {
			case "image/png":
				ext = ".png"
			case "image/jpeg":
				ext = ".jpg"
			case "image/gif":
				ext = ".gif"
			}

			// Generate a filename
			filename := fmt.Sprintf("image_%d%s", time.Now().UnixNano(), ext)

			// Write to disk
			if err := os.WriteFile(filename, imgData, 0644); err != nil {
				return fmt.Errorf("failed to write image to disk: %w", err)
			}
			fmt.Printf("[Image saved as %s]\n", filename)

		case "audio":
			dataStr, ok := c["data"].(string)
			if !ok {
				return fmt.Errorf("audio content item does not have a valid 'data' field: %v", c)
			}
			mimeType, ok := c["mimeType"].(string)
			if !ok {
				return fmt.Errorf("audio content item does not have a valid 'mimeType' field: %v", c)
			}

			// Decode base64 audio data
			audioData, err := base64.StdEncoding.DecodeString(dataStr)
			if err != nil {
				return fmt.Errorf("failed to decode base64 audio data: %w", err)
			}

			// Determine file extension from MIME type
			ext := ".audio"
			switch mimeType {
			case "audio/mpeg":
				ext = ".mp3"
			case "audio/wav":
				ext = ".wav"
			case "audio/ogg":
				ext = ".ogg"
			}

			filename := fmt.Sprintf("audio_%d%s", time.Now().UnixNano(), ext)
			if err := os.WriteFile(filename, audioData, 0644); err != nil {
				return fmt.Errorf("failed to write audio to disk: %w", err)
			}
			fmt.Printf("[Audio saved as %s]\n", filename)

		}
	}

	return nil
}
