package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/bajankristof/gerry/config"
	"github.com/bajankristof/gerry/tools"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// Configure slog to write to stderr (MCP uses stdout)
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})))

	// Load config once at startup
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Create MCP server
	s := server.NewMCPServer(
		"gerry",
		"1.0.0",
		server.WithToolCapabilities(true),
	)

	// Register tools
	tools.Inject(s, cfg)

	// Start server with stdio transport
	if err := server.ServeStdio(s); err != nil {
		log.Fatalf("Server error: %v", err)
		os.Exit(1)
	}
}
