package mcp

import (
	"context"
	"fmt"
	"net/http"

	wecombot "github.com/futuretea/go-wecom-bot"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/futuretea/wecom-bot-mcp-server/pkg/core/config"
	"github.com/futuretea/wecom-bot-mcp-server/pkg/core/logging"
	"github.com/futuretea/wecom-bot-mcp-server/pkg/core/version"
	"github.com/futuretea/wecom-bot-mcp-server/pkg/toolset"
	wecomToolset "github.com/futuretea/wecom-bot-mcp-server/pkg/toolset/wecom"
)

// Server represents the MCP server
type Server struct {
	config       *config.StaticConfig
	server       *server.MCPServer
	enabledTools []string
	bot          *wecombot.Bot
}

// NewServer creates a new MCP server with the given configuration
func NewServer(cfg *config.StaticConfig) (*Server, error) {
	serverOptions := []server.ServerOption{
		server.WithToolCapabilities(true),
		server.WithLogging(),
	}

	// Initialize WeCom bot client
	bot := wecombot.New(cfg.WeComBotKey)
	logging.Info("WeCom bot client initialized")

	s := &Server{
		config: cfg,
		server: server.NewMCPServer(version.BinaryName, version.Version, serverOptions...),
		bot:    bot,
	}

	// Register tools
	if err := s.registerTools(); err != nil {
		return nil, err
	}

	return s, nil
}

// registerTools registers all available tools based on configuration
func (s *Server) registerTools() error {
	wecomToolset := &wecomToolset.Toolset{}
	tools := wecomToolset.GetTools(s.bot)

	for _, tool := range tools {
		if !s.isToolEnabled(tool.Tool.Name) {
			continue
		}
		if err := s.registerTool(tool); err != nil {
			return fmt.Errorf("failed to register tool %s: %w", tool.Tool.Name, err)
		}
	}

	logging.Info("MCP server initialized with %d tools", len(s.enabledTools))
	return nil
}

// isToolEnabled determines if a tool should be enabled based on configuration
func (s *Server) isToolEnabled(toolName string) bool {
	// Explicitly disabled tools take highest priority
	for _, disabled := range s.config.DisabledTools {
		if disabled == toolName {
			return false
		}
	}

	// If no enabled tools specified, all non-disabled tools are enabled
	if len(s.config.EnabledTools) == 0 {
		return true
	}

	// Tool must be in the enabled list
	for _, enabled := range s.config.EnabledTools {
		if enabled == toolName {
			return true
		}
	}
	return false
}

// registerTool registers a single tool with the MCP server
func (s *Server) registerTool(tool toolset.ServerTool) error {
	handler := s.createToolHandler(tool)
	s.server.AddTool(tool.Tool, handler)
	s.enabledTools = append(s.enabledTools, tool.Tool.Name)

	logging.Info("Registered tool: %s", tool.Tool.Name)
	return nil
}

// createToolHandler creates the handler function for a tool
func (s *Server) createToolHandler(tool toolset.ServerTool) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logging.Debug("Tool %s called with params: %v", tool.Tool.Name, request.Params.Arguments)

		params := extractParams(request.Params.Arguments)
		result, err := tool.Handler(s.bot, params)
		return NewTextResult(result, err), nil
	}
}

// extractParams extracts the parameters map from the request arguments
func extractParams(args any) map[string]any {
	params, ok := args.(map[string]any)
	if !ok {
		return make(map[string]any)
	}
	return params
}

// ServeStdio starts the MCP server in stdio mode
func (s *Server) ServeStdio() error {
	logging.Info("Starting MCP server in stdio mode")
	return server.ServeStdio(s.server)
}

// ServeSse starts the MCP server in SSE mode
func (s *Server) ServeSse(baseURL string, httpServer *http.Server) *server.SSEServer {
	logging.Info("Starting MCP server in SSE mode")

	options := []server.SSEOption{
		server.WithHTTPServer(httpServer),
	}
	if baseURL != "" {
		options = append(options, server.WithBaseURL(baseURL))
	}

	return server.NewSSEServer(s.server, options...)
}

// ServeHTTP starts the MCP server in HTTP mode
func (s *Server) ServeHTTP(httpServer *http.Server) *server.StreamableHTTPServer {
	logging.Info("Starting MCP server in HTTP mode")

	options := []server.StreamableHTTPOption{
		server.WithStreamableHTTPServer(httpServer),
		server.WithStateLess(true),
	}

	return server.NewStreamableHTTPServer(s.server, options...)
}

// GetEnabledTools returns the list of enabled tools
func (s *Server) GetEnabledTools() []string {
	return s.enabledTools
}

// IsHealthy returns true if the server is properly initialized
func (s *Server) IsHealthy() bool {
	return s.bot != nil
}

// Close cleans up the server resources
func (s *Server) Close() {
	logging.Info("Closing MCP server")
}

// NewTextResult creates a standardized text result for tool responses
func NewTextResult(content string, err error) *mcp.CallToolResult {
	if err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: err.Error(),
				},
			},
		}
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: content,
			},
		},
	}
}
