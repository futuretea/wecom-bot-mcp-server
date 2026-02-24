package mcp

import (
	"errors"
	"testing"

	mcpgo "github.com/mark3labs/mcp-go/mcp"

	"github.com/futuretea/wecom-bot-mcp-server/pkg/core/config"
)

// --- isToolEnabled tests ---

func newServerWithConfig(enabled, disabled []string) *Server {
	return &Server{
		config: &config.StaticConfig{
			EnabledTools:  enabled,
			DisabledTools: disabled,
		},
	}
}

func TestIsToolEnabled_AllEnabledByDefault(t *testing.T) {
	s := newServerWithConfig(nil, nil)
	if !s.isToolEnabled("send_text") {
		t.Fatal("expected tool to be enabled by default")
	}
}

func TestIsToolEnabled_DisabledTakesPriority(t *testing.T) {
	s := newServerWithConfig([]string{"send_text"}, []string{"send_text"})
	if s.isToolEnabled("send_text") {
		t.Fatal("expected tool to be disabled when in both lists")
	}
}

func TestIsToolEnabled_OnlyEnabled(t *testing.T) {
	s := newServerWithConfig([]string{"send_text", "send_markdown"}, nil)
	if !s.isToolEnabled("send_text") {
		t.Fatal("expected send_text to be enabled")
	}
	if s.isToolEnabled("send_image") {
		t.Fatal("expected send_image to not be enabled when not in enabled list")
	}
}

func TestIsToolEnabled_OnlyDisabled(t *testing.T) {
	s := newServerWithConfig(nil, []string{"send_text"})
	if s.isToolEnabled("send_text") {
		t.Fatal("expected send_text to be disabled")
	}
	if !s.isToolEnabled("send_markdown") {
		t.Fatal("expected send_markdown to be enabled (not in disabled list)")
	}
}

// --- extractParams tests ---

func TestExtractParams_ValidMap(t *testing.T) {
	input := map[string]any{"key": "value"}
	result := extractParams(input)
	if result["key"] != "value" {
		t.Fatalf("expected 'value', got %v", result["key"])
	}
}

func TestExtractParams_NonMap(t *testing.T) {
	result := extractParams("not a map")
	if len(result) != 0 {
		t.Fatalf("expected empty map, got %v", result)
	}
}

func TestExtractParams_Nil(t *testing.T) {
	result := extractParams(nil)
	if result == nil {
		t.Fatal("expected non-nil empty map")
	}
	if len(result) != 0 {
		t.Fatalf("expected empty map, got %v", result)
	}
}

// --- NewTextResult tests ---

func TestNewTextResult_Success(t *testing.T) {
	result := NewTextResult("ok", nil)
	if result.IsError {
		t.Fatal("expected IsError to be false")
	}
	if len(result.Content) != 1 {
		t.Fatalf("expected 1 content item, got %d", len(result.Content))
	}
	tc, ok := result.Content[0].(mcpgo.TextContent)
	if !ok {
		t.Fatal("expected TextContent type")
	}
	if tc.Text != "ok" {
		t.Fatalf("expected 'ok', got %q", tc.Text)
	}
}

func TestNewTextResult_Error(t *testing.T) {
	result := NewTextResult("", errors.New("something failed"))
	if !result.IsError {
		t.Fatal("expected IsError to be true")
	}
	if len(result.Content) != 1 {
		t.Fatalf("expected 1 content item, got %d", len(result.Content))
	}
	tc, ok := result.Content[0].(mcpgo.TextContent)
	if !ok {
		t.Fatal("expected TextContent type")
	}
	if tc.Text != "something failed" {
		t.Fatalf("expected 'something failed', got %q", tc.Text)
	}
}
