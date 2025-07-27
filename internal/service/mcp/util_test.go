package mcp

import (
	"testing"
)

func TestValidateServerName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid name", "server_1", false},
		{"valid hyphen", "server-2", false},
		{"invalid slash", "server/3", true},
		{"invalid special char", "server$", true},
		{"empty", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateServerName(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateServerName(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestMergeServerToolNames(t *testing.T) {
	tests := []struct {
		server string
		tool   string
		want   string
	}{
		{"myserver", "mytool", "myserver__mytool"},
		{"myserver", "my/tool", "myserver__my/tool"},
	}
	for _, tt := range tests {
		t.Run(tt.server+"_"+tt.tool, func(t *testing.T) {
			got := mergeServerToolNames(tt.server, tt.tool)
			if got != tt.want {
				t.Errorf("mergeServerToolNames(%q, %q) = %q, want %q", tt.server, tt.tool, got, tt.want)
			}
		})
	}
}

func TestSplitServerToolName(t *testing.T) {
	tests := []struct {
		input      string
		wantServer string
		wantTool   string
		wantOK     bool
	}{
		{"server__tool", "server", "tool", true},
		{"a__b/c", "a", "b/c", true},
		{"no_separator", "", "", false},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			server, tool, ok := splitServerToolName(tt.input)
			if server != tt.wantServer || tool != tt.wantTool || ok != tt.wantOK {
				t.Errorf("splitServerToolName(%q) = (%q, %q, %v), want (%q, %q, %v)",
					tt.input, server, tool, ok, tt.wantServer, tt.wantTool, tt.wantOK)
			}
		})
	}
}

func TestIsLoopbackURL(t *testing.T) {
	tests := []struct {
		name   string
		rawURL string
		want   bool
	}{
		// IPv4 loopback
		{"IPv4 127.0.0.1", "http://127.0.0.1:8080", true},
		{"IPv4 127.0.0.1 no port", "http://127.0.0.1", true},
		{"IPv4 127.0.0.2", "http://127.0.0.2", true}, // 127.0.0.0/8 is loopback
		{"IPv4 127.255.255.255", "http://127.255.255.255", true},
		{"IPv4 0.0.0.0", "http://0.0.0.0:9000", false}, // 0.0.0.0 is not loopback, it's "any"
		// IPv6 loopback
		{"IPv6 ::1", "http://[::1]:8080", true},
		{"IPv6 ::1 no port", "http://[::1]", true},
		// Hostname loopback
		{"localhost", "http://localhost:8080", true},
		{"localhost no port", "http://localhost", true},
		{"LOCALHOST uppercase", "http://LOCALHOST", true},
		// Non-loopback IPv4
		{"IPv4 public", "http://8.8.8.8:8080", false},
		{"IPv4 private", "http://192.168.1.1", false},
		// Non-loopback IPv6
		{"IPv6 public", "http://[2001:4860:4860::8888]:443", false},
		// Hostname non-loopback
		{"example.com", "http://example.com", false},
		{"sub.domain.com", "http://sub.domain.com:1234", false},
		// Malformed URLs
		{"empty string", "", false},
		{"no scheme", "127.0.0.1:8080", false},
		{"garbage", "not a url", false},
		// Edge cases
		{"IPv4 with userinfo", "http://user:pass@127.0.0.1:8080", true},
		{"IPv6 with userinfo", "http://user:pass@[::1]:8080", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isLoopbackURL(tt.rawURL)
			if got != tt.want {
				t.Errorf("isLoopbackURL(%q) = %v, want %v", tt.rawURL, got, tt.want)
			}
		})
	}
}
