package main

import (
	"errors"
	"fmt"
	"github.com/mcpjungle/mcpjungle/cmd"
	"os"
)

func main() {
	if err := cmd.Execute(); err != nil {
		if !errors.Is(err, cmd.SilentErr) {
			_, _ = fmt.Fprintln(os.Stderr, err)
		}
		os.Exit(1)
	}
}
