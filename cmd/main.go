package main

import (
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/bajankristof/gerry/config"
	"github.com/bajankristof/gerry/tools"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// Configure slog to write to stderr (MCP uses stdout)
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})))

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	s := server.NewMCPServer(
		"gerry",
		"1.0.0",
		server.WithToolCapabilities(true),
	)
	tools.Inject(s, cfg)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

	errs := make(chan error, 1)
	go func() {
		if err := server.ServeStdio(s); err != nil {
			errs <- err
		}
	}()

	select {
	case <-sigs:
		os.Exit(0)
	case err := <-errs:
		log.Fatalf("Server error: %v", err)
	}
}
