package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/adhocteam/recreation-mcp-server/internal/config"
	"github.com/adhocteam/recreation-mcp-server/internal/mcp"
	"github.com/adhocteam/recreation-mcp-server/pkg/util"
	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	// Load .env file if it exists (ignore error if not found)
	_ = godotenv.Load()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger
	logger := util.NewLogger(cfg.LogLevel, "json")
	logger.Info("Starting Recreation MCP Server")

	// Create MCP server
	server, err := mcp.NewServer(cfg, logger)
	if err != nil {
		log.Fatalf("Failed to create MCP server: %v", err)
	}

	// Set up context with cancellation
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

	// Create stdio transport
	transport := &mcpsdk.StdioTransport{}

	// Run server
	logger.Info("MCP Server ready - communicating over stdio")
	if err := server.Run(ctx, transport); err != nil {
		logger.Error("Server error", "error", err.Error())
		cancel()
		os.Exit(1)
	}

	cancel()
	logger.Info("MCP Server shut down gracefully")
}
