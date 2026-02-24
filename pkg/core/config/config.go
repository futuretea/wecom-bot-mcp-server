package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// StaticConfig represents the static configuration for the WeCom Bot MCP Server
type StaticConfig struct {
	// Server configuration
	Port int `mapstructure:"port"`

	SSEBaseURL string `mapstructure:"sse_base_url"`

	// Logging configuration
	LogLevel int `mapstructure:"log_level"`

	// WeCom Bot configuration
	WeComBotKey string `mapstructure:"wecom_bot_key"`

	// Tool configuration
	EnabledTools  []string `mapstructure:"enabled_tools"`
	DisabledTools []string `mapstructure:"disabled_tools"`
}

// Validate validates the configuration
func (c *StaticConfig) Validate() error {
	// Validate port
	if c.Port < 0 || c.Port > 65535 {
		return fmt.Errorf("port must be between 0 and 65535, got %d", c.Port)
	}

	// Validate log level
	if c.LogLevel < 0 || c.LogLevel > 9 {
		return fmt.Errorf("log_level must be between 0 and 9, got %d", c.LogLevel)
	}

	// Validate WeCom Bot key
	if c.WeComBotKey == "" {
		return fmt.Errorf("wecom_bot_key is required")
	}

	return nil
}

// LoadConfig loads configuration from file and environment variables using Viper
// Priority: command-line flags > environment variables > config file > defaults
func LoadConfig(configPath string) (*StaticConfig, error) {
	// Use the global viper instance to access bound command-line flags
	v := viper.GetViper()

	// Set configuration file if provided
	if configPath != "" {
		v.SetConfigFile(configPath)
		v.SetConfigType("yaml")
		if err := v.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	// Configure environment variable support
	// Environment variables use WECOM_MCP_ prefix and replace - with _
	v.SetEnvPrefix("WECOM_MCP")
	v.AllowEmptyEnv(true)
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))
	v.AutomaticEnv()

	// Unmarshal configuration into struct
	config := &StaticConfig{}
	if err := v.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Validate configuration
	if err := config.Validate(); err != nil {
		return nil, err
	}

	return config, nil
}
