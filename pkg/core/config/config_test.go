package config

import (
	"strings"
	"testing"
)

func validConfig() *StaticConfig {
	return &StaticConfig{
		Port:        8080,
		LogLevel:    5,
		WeComBotKey: "test-key-123",
	}
}

func TestValidate_ValidConfig(t *testing.T) {
	cfg := validConfig()
	if err := cfg.Validate(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestValidate_StdioMode(t *testing.T) {
	cfg := validConfig()
	cfg.Port = 0 // stdio mode
	if err := cfg.Validate(); err != nil {
		t.Fatalf("expected no error for port 0, got %v", err)
	}
}

func TestValidate_PortTooHigh(t *testing.T) {
	cfg := validConfig()
	cfg.Port = 65536
	err := cfg.Validate()
	if err == nil || !strings.Contains(err.Error(), "port must be between") {
		t.Fatalf("expected port validation error, got %v", err)
	}
}

func TestValidate_PortNegative(t *testing.T) {
	cfg := validConfig()
	cfg.Port = -1
	err := cfg.Validate()
	if err == nil || !strings.Contains(err.Error(), "port must be between") {
		t.Fatalf("expected port validation error, got %v", err)
	}
}

func TestValidate_LogLevelTooHigh(t *testing.T) {
	cfg := validConfig()
	cfg.LogLevel = 10
	err := cfg.Validate()
	if err == nil || !strings.Contains(err.Error(), "log_level must be between") {
		t.Fatalf("expected log_level validation error, got %v", err)
	}
}

func TestValidate_LogLevelNegative(t *testing.T) {
	cfg := validConfig()
	cfg.LogLevel = -1
	err := cfg.Validate()
	if err == nil || !strings.Contains(err.Error(), "log_level must be between") {
		t.Fatalf("expected log_level validation error, got %v", err)
	}
}

func TestValidate_MissingWeComBotKey(t *testing.T) {
	cfg := validConfig()
	cfg.WeComBotKey = ""
	err := cfg.Validate()
	if err == nil || !strings.Contains(err.Error(), "wecom_bot_key is required") {
		t.Fatalf("expected wecom_bot_key validation error, got %v", err)
	}
}

func TestValidate_BoundaryLogLevels(t *testing.T) {
	for _, level := range []int{0, 9} {
		cfg := validConfig()
		cfg.LogLevel = level
		if err := cfg.Validate(); err != nil {
			t.Fatalf("expected no error for log_level %d, got %v", level, err)
		}
	}
}

func TestValidate_BoundaryPorts(t *testing.T) {
	for _, port := range []int{0, 65535} {
		cfg := validConfig()
		cfg.Port = port
		if err := cfg.Validate(); err != nil {
			t.Fatalf("expected no error for port %d, got %v", port, err)
		}
	}
}
