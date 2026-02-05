package mcp

import (
	"context"
	"fmt"

	"github.com/adhocteam/recreation-mcp-server/internal/api"
	"github.com/adhocteam/recreation-mcp-server/internal/cache"
	"github.com/adhocteam/recreation-mcp-server/internal/config"
	"github.com/adhocteam/recreation-mcp-server/pkg/util"
	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// Server wraps the MCP SDK server with our API clients
type Server struct {
	mcpServer        *mcpsdk.Server
	npsClient        api.NPSClient
	recreationClient api.RecreationGovClient
	weatherClient    api.WeatherClient
	logger           *util.Logger
}

// NewServer creates a new MCP server with initialized API clients
func NewServer(cfg *config.Config, logger *util.Logger) (*Server, error) {
	// Initialize cache
	cacheInstance := cache.NewCache(
		cfg.Cache.TTL,
		100, // 100MB max cache size
		cfg.Cache.Enabled,
	)

	// Start cache cleanup routine
	if cfg.Cache.Enabled {
		cacheInstance.StartCleanupRoutine(cfg.Cache.TTL / 2)
	}

	// Create HTTP client with retry logic
	httpClient := util.NewHTTPClient(
		cfg.APIs.NPS.Timeout,
		cfg.APIs.NPS.RetryAttempts,
		logger,
	)

	// Initialize API clients
	npsClient := api.NewNPSClient(
		cfg.APIs.NPS.BaseURL,
		cfg.NPSAPIKey,
		httpClient,
		cacheInstance,
		logger,
	)

	recreationClient := api.NewRecreationGovClient(
		cfg.APIs.RecreationGov.BaseURL,
		cfg.RecreationGovAPIKey,
		httpClient,
		cacheInstance,
		logger,
	)

	weatherClient := api.NewWeatherClient(
		cfg.APIs.OpenWeather.BaseURL,
		cfg.OpenWeatherAPIKey,
		httpClient,
		cacheInstance,
		logger,
	)

	// Create MCP server
	implementation := &mcpsdk.Implementation{
		Name:    "recreation-mcp-server",
		Version: "1.0.0",
	}

	mcpServer := mcpsdk.NewServer(implementation, nil)

	server := &Server{
		mcpServer:        mcpServer,
		npsClient:        npsClient,
		recreationClient: recreationClient,
		weatherClient:    weatherClient,
		logger:           logger,
	}

	// Register all tools
	if err := server.registerTools(); err != nil {
		return nil, fmt.Errorf("failed to register tools: %w", err)
	}

	logger.Info("MCP server initialized successfully")
	return server, nil
}

// Run starts the MCP server over the given transport
func (s *Server) Run(ctx context.Context, transport mcpsdk.Transport) error {
	s.logger.Info("Starting MCP server...")
	return s.mcpServer.Run(ctx, transport)
}

// registerTools registers all MCP tools with the server
func (s *Server) registerTools() error {
	// Register search_parks tool
	mcpsdk.AddTool(s.mcpServer, &mcpsdk.Tool{
		Name:        "search_parks",
		Description: "Search for national parks by name, state, or activity. Returns basic park info. For detailed planning, use a small limit (3-5 parks) then follow up with get_park_details, search_campgrounds, and get_park_alerts for each park of interest.",
	}, s.handleSearchParks)

	// Register get_park_details tool
	mcpsdk.AddTool(s.mcpServer, &mcpsdk.Tool{
		Name:        "get_park_details",
		Description: "Get detailed information about a specific national park",
	}, s.handleGetParkDetails)

	// Register get_park_alerts tool
	mcpsdk.AddTool(s.mcpServer, &mcpsdk.Tool{
		Name:        "get_park_alerts",
		Description: "Get current alerts and closures for a national park",
	}, s.handleGetParkAlerts)

	// Register search_campgrounds tool
	mcpsdk.AddTool(s.mcpServer, &mcpsdk.Tool{
		Name:        "search_campgrounds",
		Description: "Search for campgrounds in national parks. You can search by park_code (e.g. 'yose' for Yosemite), state code (e.g. 'CA'), or query text (e.g. 'yosemite'). Use get_park_details first to find a park's code if needed.",
	}, s.handleSearchCampgrounds)

	// Register search_recreation_areas tool
	mcpsdk.AddTool(s.mcpServer, &mcpsdk.Tool{
		Name:        "search_recreation_areas",
		Description: "Search recreation areas on Recreation.gov",
	}, s.handleSearchRecreationAreas)

	// Register get_facility_details tool
	mcpsdk.AddTool(s.mcpServer, &mcpsdk.Tool{
		Name:        "get_facility_details",
		Description: "Get detailed information about a specific Recreation.gov facility",
	}, s.handleGetFacilityDetails)

	// Register get_weather tool
	mcpsdk.AddTool(s.mcpServer, &mcpsdk.Tool{
		Name:        "get_weather",
		Description: "Get current weather conditions for a location",
	}, s.handleGetWeather)

	// Register get_weather_forecast tool
	mcpsdk.AddTool(s.mcpServer, &mcpsdk.Tool{
		Name:        "get_weather_forecast",
		Description: "Get 5-day weather forecast for a location",
	}, s.handleGetWeatherForecast)

	// Register list_activities tool
	mcpsdk.AddTool(s.mcpServer, &mcpsdk.Tool{
		Name:        "list_activities",
		Description: "List all available activities across NPS and Recreation.gov",
	}, s.handleListActivities)

	s.logger.Info("Registered 9 MCP tools")
	return nil
}
