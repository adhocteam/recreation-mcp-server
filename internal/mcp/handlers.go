package mcp

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/markheadd/recreation-mcp-server/internal/api"
	"github.com/markheadd/recreation-mcp-server/internal/models"
	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// Tool input and output types

// SearchParksInput defines the input for search_parks tool
type SearchParksInput struct {
	Query    string `json:"query,omitempty" jsonschema:"Search text to filter parks by name or description"`
	State    string `json:"state,omitempty" jsonschema:"Two-letter state code (e.g., CA, CO, UT)"`
	Activity string `json:"activity,omitempty" jsonschema:"Activity name to filter parks (e.g., Hiking, Camping)"`
	Limit    int    `json:"limit,omitempty" jsonschema:"Maximum number of results to return. Use 3-5 for detailed planning where you'll call other tools for each park. Default=10"`
}

// SearchParksOutput defines the output for search_parks tool
type SearchParksOutput struct {
	Parks []models.Park `json:"parks" jsonschema:"List of matching national parks"`
	Count int           `json:"count" jsonschema:"Number of parks returned"`
}

// Tool handlers

// handleSearchParks handles the search_parks tool invocation
func (s *Server) handleSearchParks(ctx context.Context, req *mcpsdk.CallToolRequest, input SearchParksInput) (*mcpsdk.CallToolResult, SearchParksOutput, error) {
	s.logger.Info("Handling search_parks request", "query", input.Query, "state", input.State)

	// Set default limit if not specified
	if input.Limit == 0 {
		input.Limit = 10
	}

	// Build search parameters
	params := api.SearchParksParams{
		Query:    input.Query,
		State:    input.State,
		Activity: input.Activity,
		Limit:    input.Limit,
		Start:    0,
	}

	// Execute search
	parks, err := s.npsClient.SearchParks(ctx, params)
	if err != nil {
		s.logger.Error("Failed to search parks", "error", err.Error())
		return mcpError("api_error", fmt.Sprintf("Failed to search parks: %v", err)), SearchParksOutput{Parks: []models.Park{}, Count: 0}, nil
	}

	// Ensure parks is never nil
	if parks == nil {
		parks = []models.Park{}
	}

	output := SearchParksOutput{
		Parks: parks,
		Count: len(parks),
	}

	return nil, output, nil
}

// GetParkDetailsInput defines the input for get_park_details tool
type GetParkDetailsInput struct {
	ParkCode string `json:"park_code" jsonschema:"Four-letter park code (e.g., 'yose' for Yosemite, 'grca' for Grand Canyon)"`
}

// handleGetParkDetails handles the get_park_details tool invocation
func (s *Server) handleGetParkDetails(ctx context.Context, req *mcpsdk.CallToolRequest, input GetParkDetailsInput) (*mcpsdk.CallToolResult, *models.Park, error) {
	s.logger.Info("Handling get_park_details request", "parkCode", input.ParkCode)

	if input.ParkCode == "" {
		return mcpError("invalid_params", "park_code is required"), nil, nil
	}

	park, err := s.npsClient.GetParkDetails(ctx, input.ParkCode)
	if err != nil {
		s.logger.Error("Failed to get park details", "error", err.Error())
		return mcpError("api_error", fmt.Sprintf("Failed to get park details: %v", err)), nil, nil
	}

	return nil, park, nil
}

// GetParkAlertsInput defines the input for get_park_alerts tool
type GetParkAlertsInput struct {
	ParkCode string `json:"park_code" jsonschema:"Four-letter park code (e.g., 'yose' for Yosemite)"`
}

// GetParkAlertsOutput defines the output for get_park_alerts tool
type GetParkAlertsOutput struct {
	Alerts []models.Alert `json:"alerts" jsonschema:"List of current alerts and closures"`
	Count  int            `json:"count" jsonschema:"Number of alerts"`
}

// handleGetParkAlerts handles the get_park_alerts tool invocation
func (s *Server) handleGetParkAlerts(ctx context.Context, req *mcpsdk.CallToolRequest, input GetParkAlertsInput) (*mcpsdk.CallToolResult, GetParkAlertsOutput, error) {
	s.logger.Info("Handling get_park_alerts request", "parkCode", input.ParkCode)

	if input.ParkCode == "" {
		return mcpError("invalid_params", "park_code is required"), GetParkAlertsOutput{}, nil
	}

	alerts, err := s.npsClient.GetAlerts(ctx, input.ParkCode)
	if err != nil {
		s.logger.Error("Failed to get park alerts", "error", err.Error())
		return mcpError("api_error", fmt.Sprintf("Failed to get alerts: %v", err)), GetParkAlertsOutput{Alerts: []models.Alert{}, Count: 0}, nil
	}

	// Ensure alerts is never nil
	if alerts == nil {
		alerts = []models.Alert{}
	}

	output := GetParkAlertsOutput{
		Alerts: alerts,
		Count:  len(alerts),
	}

	return nil, output, nil
}

// SearchCampgroundsInput defines the input for search_campgrounds tool
type SearchCampgroundsInput struct {
	ParkCode string `json:"park_code,omitempty" jsonschema:"NPS park code (e.g. 'yose' for Yosemite) - get this from search_parks or get_park_details first"`
	State    string `json:"state,omitempty" jsonschema:"Two-letter state code (e.g. 'CA', 'CO')"`
	Query    string `json:"query,omitempty" jsonschema:"Search text for campground or park names (e.g. 'yosemite', 'yellowstone')"`
	Limit    int    `json:"limit,omitempty" jsonschema:"Maximum number of results,default=10"`
}

// SearchCampgroundsOutput defines the output for search_campgrounds tool
type SearchCampgroundsOutput struct {
	Campgrounds []models.Campground `json:"campgrounds" jsonschema:"List of matching campgrounds"`
	Count       int                 `json:"count" jsonschema:"Number of campgrounds returned"`
}

// handleSearchCampgrounds handles the search_campgrounds tool invocation
func (s *Server) handleSearchCampgrounds(ctx context.Context, req *mcpsdk.CallToolRequest, input SearchCampgroundsInput) (*mcpsdk.CallToolResult, SearchCampgroundsOutput, error) {
	s.logger.Info("Handling search_campgrounds request", "parkCode", input.ParkCode, "state", input.State)

	if input.Limit == 0 {
		input.Limit = 10
	}

	params := api.SearchCampgroundsParams{
		ParkCode: input.ParkCode,
		State:    input.State,
		Query:    input.Query,
		Limit:    input.Limit,
		Start:    0,
	}

	campgrounds, err := s.npsClient.GetCampgrounds(ctx, params)
	if err != nil {
		s.logger.Error("Failed to search campgrounds", "error", err.Error())
		return mcpError("api_error", fmt.Sprintf("Failed to search campgrounds: %v", err)), SearchCampgroundsOutput{Campgrounds: []models.Campground{}, Count: 0}, nil
	}

	// Ensure campgrounds is never nil
	if campgrounds == nil {
		campgrounds = []models.Campground{}
	}

	output := SearchCampgroundsOutput{
		Campgrounds: campgrounds,
		Count:       len(campgrounds),
	}

	return nil, output, nil
}

// SearchRecreationAreasInput defines the input for search_recreation_areas tool
type SearchRecreationAreasInput struct {
	Query    string `json:"query,omitempty" jsonschema:"Search text for recreation area names"`
	State    string `json:"state,omitempty" jsonschema:"Two-letter state code"`
	Activity string `json:"activity,omitempty" jsonschema:"Activity name to filter areas"`
	Limit    int    `json:"limit,omitempty" jsonschema:"Maximum number of results,default=10"`
}

// SearchRecreationAreasOutput defines the output for search_recreation_areas tool
type SearchRecreationAreasOutput struct {
	Areas []models.RecreationArea `json:"areas" jsonschema:"List of matching recreation areas"`
	Count int                     `json:"count" jsonschema:"Number of areas returned"`
}

// handleSearchRecreationAreas handles the search_recreation_areas tool invocation
func (s *Server) handleSearchRecreationAreas(ctx context.Context, req *mcpsdk.CallToolRequest, input SearchRecreationAreasInput) (*mcpsdk.CallToolResult, SearchRecreationAreasOutput, error) {
	s.logger.Info("Handling search_recreation_areas request", "query", input.Query, "state", input.State)

	if input.Limit == 0 {
		input.Limit = 10
	}

	params := api.SearchRecAreasParams{
		Query:    input.Query,
		State:    input.State,
		Activity: input.Activity,
		Limit:    input.Limit,
		Offset:   0,
	}

	areas, err := s.recreationClient.SearchRecreationAreas(ctx, params)
	if err != nil {
		s.logger.Error("Failed to search recreation areas", "error", err.Error())
		return mcpError("api_error", fmt.Sprintf("Failed to search recreation areas: %v", err)), SearchRecreationAreasOutput{Areas: []models.RecreationArea{}, Count: 0}, nil
	}

	// Ensure areas is never nil
	if areas == nil {
		areas = []models.RecreationArea{}
	}

	output := SearchRecreationAreasOutput{
		Areas: areas,
		Count: len(areas),
	}

	return nil, output, nil
}

// GetFacilityDetailsInput defines the input for get_facility_details tool
type GetFacilityDetailsInput struct {
	FacilityID string `json:"facility_id" jsonschema:"Recreation.gov facility ID"`
}

// handleGetFacilityDetails handles the get_facility_details tool invocation
func (s *Server) handleGetFacilityDetails(ctx context.Context, req *mcpsdk.CallToolRequest, input GetFacilityDetailsInput) (*mcpsdk.CallToolResult, *models.Facility, error) {
	s.logger.Info("Handling get_facility_details request", "facilityID", input.FacilityID)

	if input.FacilityID == "" {
		return mcpError("invalid_params", "facility_id is required"), nil, nil
	}

	facility, err := s.recreationClient.GetFacilityDetails(ctx, input.FacilityID)
	if err != nil {
		s.logger.Error("Failed to get facility details", "error", err.Error())
		return mcpError("api_error", fmt.Sprintf("Failed to get facility details: %v", err)), nil, nil
	}

	return nil, facility, nil
}

// GetWeatherInput defines the input for get_weather tool
type GetWeatherInput struct {
	Latitude  float64 `json:"latitude" jsonschema:"Latitude of the location"`
	Longitude float64 `json:"longitude" jsonschema:"Longitude of the location"`
	Units     string  `json:"units,omitempty" jsonschema:"Temperature units: 'metric' or 'imperial',default=imperial"`
}

// handleGetWeather handles the get_weather tool invocation
func (s *Server) handleGetWeather(ctx context.Context, req *mcpsdk.CallToolRequest, input GetWeatherInput) (*mcpsdk.CallToolResult, *models.Weather, error) {
	s.logger.Info("Handling get_weather request", "lat", input.Latitude, "lon", input.Longitude)

	if input.Units == "" {
		input.Units = "imperial"
	}

	weather, err := s.weatherClient.GetCurrentWeather(ctx, input.Latitude, input.Longitude, input.Units)
	if err != nil {
		s.logger.Error("Failed to get weather", "error", err.Error())
		return mcpError("api_error", fmt.Sprintf("Failed to get weather: %v", err)), nil, nil
	}

	return nil, weather, nil
}

// GetWeatherForecastInput defines the input for get_weather_forecast tool
type GetWeatherForecastInput struct {
	Latitude  float64 `json:"latitude" jsonschema:"Latitude of the location"`
	Longitude float64 `json:"longitude" jsonschema:"Longitude of the location"`
	Units     string  `json:"units,omitempty" jsonschema:"Temperature units: 'metric' or 'imperial',default=imperial"`
}

// GetWeatherForecastOutput defines the output for get_weather_forecast tool
type GetWeatherForecastOutput struct {
	Forecasts []models.WeatherForecast `json:"forecasts" jsonschema:"5-day weather forecast"`
	Count     int                      `json:"count" jsonschema:"Number of forecast periods"`
}

// handleGetWeatherForecast handles the get_weather_forecast tool invocation
func (s *Server) handleGetWeatherForecast(ctx context.Context, req *mcpsdk.CallToolRequest, input GetWeatherForecastInput) (*mcpsdk.CallToolResult, GetWeatherForecastOutput, error) {
	s.logger.Info("Handling get_weather_forecast request", "lat", input.Latitude, "lon", input.Longitude)

	if input.Units == "" {
		input.Units = "imperial"
	}

	forecasts, err := s.weatherClient.GetWeatherForecast(ctx, input.Latitude, input.Longitude, input.Units)
	if err != nil {
		s.logger.Error("Failed to get weather forecast", "error", err.Error())
		return mcpError("api_error", fmt.Sprintf("Failed to get weather forecast: %v", err)), GetWeatherForecastOutput{Forecasts: []models.WeatherForecast{}, Count: 0}, nil
	}

	// Ensure forecasts is never nil
	if forecasts == nil {
		forecasts = []models.WeatherForecast{}
	}

	output := GetWeatherForecastOutput{
		Forecasts: forecasts,
		Count:     len(forecasts),
	}

	return nil, output, nil
}

// ListActivitiesInput defines the input for list_activities tool
type ListActivitiesInput struct {
	Source string `json:"source,omitempty" jsonschema:"Data source: 'nps', 'recreation_gov', or 'all',default=all"`
	Limit  int    `json:"limit,omitempty" jsonschema:"Maximum number of results,default=50"`
}

// ListActivitiesOutput defines the output for list_activities tool
type ListActivitiesOutput struct {
	NPSActivities    []models.Activity           `json:"nps_activities,omitempty" jsonschema:"Activities from National Park Service"`
	RecGovActivities []models.RecreationActivity `json:"recgov_activities,omitempty" jsonschema:"Activities from Recreation.gov"`
	Count            int                         `json:"count" jsonschema:"Total number of activities returned"`
}

// handleListActivities handles the list_activities tool invocation
func (s *Server) handleListActivities(ctx context.Context, req *mcpsdk.CallToolRequest, input ListActivitiesInput) (*mcpsdk.CallToolResult, ListActivitiesOutput, error) {
	s.logger.Info("Handling list_activities request", "source", input.Source)

	if input.Source == "" {
		input.Source = "all"
	}
	if input.Limit == 0 {
		input.Limit = 50
	}

	// Initialize output with empty arrays
	output := ListActivitiesOutput{
		NPSActivities:    []models.Activity{},
		RecGovActivities: []models.RecreationActivity{},
	}

	// Fetch NPS activities
	if input.Source == "nps" || input.Source == "all" {
		npsActivities, err := s.npsClient.GetActivities(ctx)
		if err != nil {
			s.logger.Error("Failed to get NPS activities", "error", err.Error())
			return mcpError("api_error", fmt.Sprintf("Failed to get NPS activities: %v", err)), ListActivitiesOutput{NPSActivities: []models.Activity{}, RecGovActivities: []models.RecreationActivity{}, Count: 0}, nil
		}
		if npsActivities == nil {
			npsActivities = []models.Activity{}
		}
		output.NPSActivities = npsActivities
	}

	// Fetch Recreation.gov activities
	if input.Source == "recreation_gov" || input.Source == "all" {
		recGovActivities, err := s.recreationClient.GetActivities(ctx)
		if err != nil {
			s.logger.Error("Failed to get Recreation.gov activities", "error", err.Error())
			return mcpError("api_error", fmt.Sprintf("Failed to get Recreation.gov activities: %v", err)), ListActivitiesOutput{NPSActivities: []models.Activity{}, RecGovActivities: []models.RecreationActivity{}, Count: 0}, nil
		}
		if recGovActivities == nil {
			recGovActivities = []models.RecreationActivity{}
		}
		output.RecGovActivities = recGovActivities
	}

	output.Count = len(output.NPSActivities) + len(output.RecGovActivities)

	return nil, output, nil
}

// mcpError creates an MCP-compliant error response
func mcpError(code, message string) *mcpsdk.CallToolResult {
	errorContent := map[string]string{
		"error":   code,
		"message": message,
	}

	errorJSON, _ := json.Marshal(errorContent)

	return &mcpsdk.CallToolResult{
		Content: []mcpsdk.Content{
			&mcpsdk.TextContent{Text: string(errorJSON)},
		},
		IsError: true,
	}
}
