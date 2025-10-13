package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/markheadd/recreation-mcp-server/internal/cache"
	"github.com/markheadd/recreation-mcp-server/internal/models"
	"github.com/markheadd/recreation-mcp-server/pkg/util"
)

// NPSClient defines the interface for the National Park Service API client
type NPSClient interface {
	SearchParks(ctx context.Context, params SearchParksParams) ([]models.Park, error)
	GetParkDetails(ctx context.Context, parkCode string) (*models.Park, error)
	GetAlerts(ctx context.Context, parkCode string) ([]models.Alert, error)
	GetCampgrounds(ctx context.Context, params SearchCampgroundsParams) ([]models.Campground, error)
	GetActivities(ctx context.Context) ([]models.Activity, error)
	GetVisitorCenters(ctx context.Context, parkCode string) ([]interface{}, error)
}

// SearchParksParams holds parameters for searching parks
type SearchParksParams struct {
	Query    string
	State    string
	Activity string
	Limit    int
	Start    int
}

// SearchCampgroundsParams holds parameters for searching campgrounds
type SearchCampgroundsParams struct {
	ParkCode string
	State    string
	Query    string
	Limit    int
	Start    int
}

// npsClient implements the NPSClient interface
type npsClient struct {
	baseURL    string
	apiKey     string
	httpClient *util.HTTPClient
	cache      *cache.Cache
	logger     *util.Logger
}

// NewNPSClient creates a new National Park Service API client
func NewNPSClient(baseURL, apiKey string, httpClient *util.HTTPClient, cache *cache.Cache, logger *util.Logger) NPSClient {
	return &npsClient{
		baseURL:    baseURL,
		apiKey:     apiKey,
		httpClient: httpClient,
		cache:      cache,
		logger:     logger,
	}
}

// SearchParks searches for national parks
func (c *npsClient) SearchParks(ctx context.Context, params SearchParksParams) ([]models.Park, error) {
	c.logger.Debug("Searching parks", "query", params.Query, "state", params.State, "activity", params.Activity)

	endpoint := fmt.Sprintf("%s/parks", c.baseURL)
	queryParams := c.buildParksQueryParams(params)
	fullURL := fmt.Sprintf("%s?%s", endpoint, queryParams.Encode())

	cacheKey := cache.GenerateKey("nps:parks", map[string]string{
		"query":    params.Query,
		"state":    params.State,
		"activity": params.Activity,
		"limit":    strconv.Itoa(params.Limit),
		"start":    strconv.Itoa(params.Start),
	})

	if cached, found := c.cache.Get(cacheKey); found {
		c.logger.Debug("Cache hit for parks search")
		var parks []models.Park
		if err := json.Unmarshal(cached, &parks); err == nil {
			return parks, nil
		}
	}

	resp, err := c.httpClient.Get(ctx, fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to search parks: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := util.ReadResponseBody(resp)
	if err != nil {
		return nil, err
	}

	var npsResp models.NPSResponse
	if err := json.Unmarshal(body, &npsResp); err != nil {
		return nil, fmt.Errorf("failed to parse parks response: %w", err)
	}

	parksData, err := json.Marshal(npsResp.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal parks data: %w", err)
	}

	var parks []models.Park
	if err := json.Unmarshal(parksData, &parks); err != nil {
		return nil, fmt.Errorf("failed to unmarshal parks: %w", err)
	}

	if cacheData, err := json.Marshal(parks); err == nil {
		_ = c.cache.Set(cacheKey, cacheData)
	}

	c.logger.Info("Parks search completed", "count", len(parks))
	return parks, nil
}

// GetParkDetails retrieves detailed information about a specific park
func (c *npsClient) GetParkDetails(ctx context.Context, parkCode string) (*models.Park, error) {
	c.logger.Debug("Getting park details", "parkCode", parkCode)

	endpoint := fmt.Sprintf("%s/parks", c.baseURL)
	queryParams := url.Values{}
	queryParams.Add("parkCode", parkCode)
	queryParams.Add("api_key", c.apiKey)
	fullURL := fmt.Sprintf("%s?%s", endpoint, queryParams.Encode())

	cacheKey := cache.GenerateKey("nps:park", map[string]string{"parkCode": parkCode})

	if cached, found := c.cache.Get(cacheKey); found {
		c.logger.Debug("Cache hit for park details")
		var park models.Park
		if err := json.Unmarshal(cached, &park); err == nil {
			return &park, nil
		}
	}

	resp, err := c.httpClient.Get(ctx, fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get park details: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := util.ReadResponseBody(resp)
	if err != nil {
		return nil, err
	}

	var npsResp models.NPSResponse
	if err := json.Unmarshal(body, &npsResp); err != nil {
		return nil, fmt.Errorf("failed to parse park response: %w", err)
	}

	parksData, err := json.Marshal(npsResp.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal park data: %w", err)
	}

	var parks []models.Park
	if err := json.Unmarshal(parksData, &parks); err != nil {
		return nil, fmt.Errorf("failed to unmarshal park: %w", err)
	}

	if len(parks) == 0 {
		return nil, fmt.Errorf("park not found: %s", parkCode)
	}

	park := &parks[0]

	if cacheData, err := json.Marshal(park); err == nil {
		_ = c.cache.Set(cacheKey, cacheData)
	}

	c.logger.Info("Park details retrieved", "parkCode", parkCode)
	return park, nil
}

// GetAlerts retrieves alerts for a specific park
func (c *npsClient) GetAlerts(ctx context.Context, parkCode string) ([]models.Alert, error) {
	c.logger.Debug("Getting alerts", "parkCode", parkCode)

	endpoint := fmt.Sprintf("%s/alerts", c.baseURL)
	queryParams := url.Values{}
	queryParams.Add("parkCode", parkCode)
	queryParams.Add("api_key", c.apiKey)
	fullURL := fmt.Sprintf("%s?%s", endpoint, queryParams.Encode())

	cacheKey := cache.GenerateKey("nps:alerts", map[string]string{"parkCode": parkCode})

	if cached, found := c.cache.Get(cacheKey); found {
		c.logger.Debug("Cache hit for alerts")
		var alerts []models.Alert
		if err := json.Unmarshal(cached, &alerts); err == nil {
			return alerts, nil
		}
	}

	resp, err := c.httpClient.Get(ctx, fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get alerts: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := util.ReadResponseBody(resp)
	if err != nil {
		return nil, err
	}

	var npsResp models.NPSResponse
	if err := json.Unmarshal(body, &npsResp); err != nil {
		return nil, fmt.Errorf("failed to parse alerts response: %w", err)
	}

	alertsData, err := json.Marshal(npsResp.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal alerts data: %w", err)
	}

	var alerts []models.Alert
	if err := json.Unmarshal(alertsData, &alerts); err != nil {
		return nil, fmt.Errorf("failed to unmarshal alerts: %w", err)
	}

	if cacheData, err := json.Marshal(alerts); err == nil {
		_ = c.cache.Set(cacheKey, cacheData)
	}

	c.logger.Info("Alerts retrieved", "parkCode", parkCode, "count", len(alerts))
	return alerts, nil
}

// GetCampgrounds searches for campgrounds
func (c *npsClient) GetCampgrounds(ctx context.Context, params SearchCampgroundsParams) ([]models.Campground, error) {
	c.logger.Debug("Searching campgrounds", "parkCode", params.ParkCode, "state", params.State, "query", params.Query)

	endpoint := fmt.Sprintf("%s/campgrounds", c.baseURL)
	queryParams := c.buildCampgroundsQueryParams(params)
	fullURL := fmt.Sprintf("%s?%s", endpoint, queryParams.Encode())

	cacheKey := cache.GenerateKey("nps:campgrounds", map[string]string{
		"parkCode": params.ParkCode,
		"state":    params.State,
		"query":    params.Query,
		"limit":    strconv.Itoa(params.Limit),
		"start":    strconv.Itoa(params.Start),
	})

	if cached, found := c.cache.Get(cacheKey); found {
		c.logger.Debug("Cache hit for campgrounds search")
		var campgrounds []models.Campground
		if err := json.Unmarshal(cached, &campgrounds); err == nil {
			return campgrounds, nil
		}
	}

	resp, err := c.httpClient.Get(ctx, fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to search campgrounds: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := util.ReadResponseBody(resp)
	if err != nil {
		return nil, err
	}

	var npsResp models.NPSResponse
	if err := json.Unmarshal(body, &npsResp); err != nil {
		return nil, fmt.Errorf("failed to parse campgrounds response: %w", err)
	}

	campgroundsData, err := json.Marshal(npsResp.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal campgrounds data: %w", err)
	}

	var campgrounds []models.Campground
	if err := json.Unmarshal(campgroundsData, &campgrounds); err != nil {
		return nil, fmt.Errorf("failed to unmarshal campgrounds: %w", err)
	}

	if cacheData, err := json.Marshal(campgrounds); err == nil {
		_ = c.cache.Set(cacheKey, cacheData)
	}

	c.logger.Info("Campgrounds search completed", "count", len(campgrounds))
	return campgrounds, nil
}

// GetActivities retrieves all available activities
func (c *npsClient) GetActivities(ctx context.Context) ([]models.Activity, error) {
	c.logger.Debug("Getting activities")

	endpoint := fmt.Sprintf("%s/activities", c.baseURL)
	queryParams := url.Values{}
	queryParams.Add("api_key", c.apiKey)
	fullURL := fmt.Sprintf("%s?%s", endpoint, queryParams.Encode())

	cacheKey := cache.GenerateKey("nps:activities", map[string]string{})

	if cached, found := c.cache.Get(cacheKey); found {
		c.logger.Debug("Cache hit for activities")
		var activities []models.Activity
		if err := json.Unmarshal(cached, &activities); err == nil {
			return activities, nil
		}
	}

	resp, err := c.httpClient.Get(ctx, fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get activities: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := util.ReadResponseBody(resp)
	if err != nil {
		return nil, err
	}

	var npsResp models.NPSResponse
	if err := json.Unmarshal(body, &npsResp); err != nil {
		return nil, fmt.Errorf("failed to parse activities response: %w", err)
	}

	activitiesData, err := json.Marshal(npsResp.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal activities data: %w", err)
	}

	var activities []models.Activity
	if err := json.Unmarshal(activitiesData, &activities); err != nil {
		return nil, fmt.Errorf("failed to unmarshal activities: %w", err)
	}

	if cacheData, err := json.Marshal(activities); err == nil {
		_ = c.cache.Set(cacheKey, cacheData)
	}

	c.logger.Info("Activities retrieved", "count", len(activities))
	return activities, nil
}

// GetVisitorCenters retrieves visitor centers for a park
func (c *npsClient) GetVisitorCenters(ctx context.Context, parkCode string) ([]interface{}, error) {
	c.logger.Debug("Getting visitor centers", "parkCode", parkCode)

	endpoint := fmt.Sprintf("%s/visitorcenters", c.baseURL)
	queryParams := url.Values{}
	queryParams.Add("parkCode", parkCode)
	queryParams.Add("api_key", c.apiKey)
	fullURL := fmt.Sprintf("%s?%s", endpoint, queryParams.Encode())

	cacheKey := cache.GenerateKey("nps:visitorcenters", map[string]string{"parkCode": parkCode})

	if cached, found := c.cache.Get(cacheKey); found {
		c.logger.Debug("Cache hit for visitor centers")
		var centers []interface{}
		if err := json.Unmarshal(cached, &centers); err == nil {
			return centers, nil
		}
	}

	resp, err := c.httpClient.Get(ctx, fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get visitor centers: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := util.ReadResponseBody(resp)
	if err != nil {
		return nil, err
	}

	var npsResp models.NPSResponse
	if err := json.Unmarshal(body, &npsResp); err != nil {
		return nil, fmt.Errorf("failed to parse visitor centers response: %w", err)
	}

	centersData, err := json.Marshal(npsResp.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal visitor centers data: %w", err)
	}

	var centers []interface{}
	if err := json.Unmarshal(centersData, &centers); err != nil {
		return nil, fmt.Errorf("failed to unmarshal visitor centers: %w", err)
	}

	if cacheData, err := json.Marshal(centers); err == nil {
		_ = c.cache.Set(cacheKey, cacheData)
	}

	c.logger.Info("Visitor centers retrieved", "parkCode", parkCode, "count", len(centers))
	return centers, nil
}

// buildParksQueryParams builds URL query parameters for parks search
func (c *npsClient) buildParksQueryParams(params SearchParksParams) url.Values {
	queryParams := url.Values{}
	queryParams.Add("api_key", c.apiKey)

	if params.Query != "" {
		queryParams.Add("q", params.Query)
	}
	if params.State != "" {
		queryParams.Add("stateCode", params.State)
	}
	if params.Limit > 0 {
		queryParams.Add("limit", strconv.Itoa(params.Limit))
	}
	if params.Start > 0 {
		queryParams.Add("start", strconv.Itoa(params.Start))
	}

	return queryParams
}

// buildCampgroundsQueryParams builds URL query parameters for campgrounds search
func (c *npsClient) buildCampgroundsQueryParams(params SearchCampgroundsParams) url.Values {
	queryParams := url.Values{}
	queryParams.Add("api_key", c.apiKey)

	if params.ParkCode != "" {
		queryParams.Add("parkCode", params.ParkCode)
	}
	if params.State != "" {
		queryParams.Add("stateCode", params.State)
	}
	if params.Query != "" {
		queryParams.Add("q", params.Query)
	}
	if params.Limit > 0 {
		queryParams.Add("limit", strconv.Itoa(params.Limit))
	}
	if params.Start > 0 {
		queryParams.Add("start", strconv.Itoa(params.Start))
	}

	return queryParams
}
