package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/adhocteam/recreation-mcp-server/internal/cache"
	"github.com/adhocteam/recreation-mcp-server/internal/models"
	"github.com/adhocteam/recreation-mcp-server/pkg/util"
)

// RecreationGovClient defines the interface for the Recreation.gov API client
type RecreationGovClient interface {
	SearchRecreationAreas(ctx context.Context, params SearchRecAreasParams) ([]models.RecreationArea, error)
	GetRecreationAreaDetails(ctx context.Context, recAreaID string) (*models.RecreationArea, error)
	SearchFacilities(ctx context.Context, params SearchFacilitiesParams) ([]models.Facility, error)
	GetFacilityDetails(ctx context.Context, facilityID string) (*models.Facility, error)
	GetCampsites(ctx context.Context, params SearchCampsitesParams) ([]models.Campsite, error)
	GetActivities(ctx context.Context) ([]models.RecreationActivity, error)
}

// SearchRecAreasParams holds parameters for searching recreation areas
type SearchRecAreasParams struct {
	Query    string
	State    string
	Activity string
	Limit    int
	Offset   int
}

// SearchFacilitiesParams holds parameters for searching facilities
type SearchFacilitiesParams struct {
	Query    string
	State    string
	Activity string
	Limit    int
	Offset   int
}

// SearchCampsitesParams holds parameters for searching campsites
type SearchCampsitesParams struct {
	Query  string
	State  string
	Limit  int
	Offset int
}

// recreationGovClient implements the RecreationGovClient interface
type recreationGovClient struct {
	baseURL    string
	apiKey     string
	httpClient *util.HTTPClient
	cache      *cache.Cache
	logger     *util.Logger
}

// NewRecreationGovClient creates a new Recreation.gov API client
func NewRecreationGovClient(baseURL, apiKey string, httpClient *util.HTTPClient, cache *cache.Cache, logger *util.Logger) RecreationGovClient {
	return &recreationGovClient{
		baseURL:    baseURL,
		apiKey:     apiKey,
		httpClient: httpClient,
		cache:      cache,
		logger:     logger,
	}
}

// SearchRecreationAreas searches for recreation areas
func (c *recreationGovClient) SearchRecreationAreas(ctx context.Context, params SearchRecAreasParams) ([]models.RecreationArea, error) {
	c.logger.Debug("Searching recreation areas", "query", params.Query, "state", params.State)

	endpoint := fmt.Sprintf("%s/recareas", c.baseURL)
	queryParams := c.buildRecAreasQueryParams(params)
	fullURL := fmt.Sprintf("%s?%s", endpoint, queryParams.Encode())

	cacheKey := cache.GenerateKey("recgov:recareas", map[string]string{
		"query":    params.Query,
		"state":    params.State,
		"activity": params.Activity,
		"limit":    strconv.Itoa(params.Limit),
		"offset":   strconv.Itoa(params.Offset),
	})

	if cached, found := c.cache.Get(cacheKey); found {
		c.logger.Debug("Cache hit for recreation areas search")
		var areas []models.RecreationArea
		if err := json.Unmarshal(cached, &areas); err == nil {
			return areas, nil
		}
	}

	headers := map[string]string{
		"apikey": c.apiKey,
	}

	resp, err := c.httpClient.Get(ctx, fullURL, headers)
	if err != nil {
		return nil, fmt.Errorf("failed to search recreation areas: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := util.ReadResponseBody(resp)
	if err != nil {
		return nil, err
	}

	var recgovResp models.RecreationGovResponse
	if err := json.Unmarshal(body, &recgovResp); err != nil {
		return nil, fmt.Errorf("failed to parse recreation areas response: %w", err)
	}

	areasData, err := json.Marshal(recgovResp.RECDATA)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal recreation areas data: %w", err)
	}

	var areas []models.RecreationArea
	if err := json.Unmarshal(areasData, &areas); err != nil {
		return nil, fmt.Errorf("failed to unmarshal recreation areas: %w", err)
	}

	if cacheData, err := json.Marshal(areas); err == nil {
		_ = c.cache.Set(cacheKey, cacheData)
	}

	c.logger.Info("Recreation areas search completed", "count", len(areas))
	return areas, nil
}

// GetRecreationAreaDetails retrieves detailed information about a specific recreation area
func (c *recreationGovClient) GetRecreationAreaDetails(ctx context.Context, recAreaID string) (*models.RecreationArea, error) {
	c.logger.Debug("Getting recreation area details", "recAreaID", recAreaID)

	endpoint := fmt.Sprintf("%s/recareas/%s", c.baseURL, recAreaID)

	cacheKey := cache.GenerateKey("recgov:recarea", map[string]string{"recAreaID": recAreaID})

	if cached, found := c.cache.Get(cacheKey); found {
		c.logger.Debug("Cache hit for recreation area details")
		var area models.RecreationArea
		if err := json.Unmarshal(cached, &area); err == nil {
			return &area, nil
		}
	}

	headers := map[string]string{
		"apikey": c.apiKey,
	}

	resp, err := c.httpClient.Get(ctx, endpoint, headers)
	if err != nil {
		return nil, fmt.Errorf("failed to get recreation area details: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := util.ReadResponseBody(resp)
	if err != nil {
		return nil, err
	}

	var area models.RecreationArea
	if err := json.Unmarshal(body, &area); err != nil {
		return nil, fmt.Errorf("failed to unmarshal recreation area: %w", err)
	}

	if cacheData, err := json.Marshal(area); err == nil {
		_ = c.cache.Set(cacheKey, cacheData)
	}

	c.logger.Info("Recreation area details retrieved", "recAreaID", recAreaID)
	return &area, nil
}

// SearchFacilities searches for facilities
func (c *recreationGovClient) SearchFacilities(ctx context.Context, params SearchFacilitiesParams) ([]models.Facility, error) {
	c.logger.Debug("Searching facilities", "query", params.Query, "state", params.State)

	endpoint := fmt.Sprintf("%s/facilities", c.baseURL)
	queryParams := c.buildFacilitiesQueryParams(params)
	fullURL := fmt.Sprintf("%s?%s", endpoint, queryParams.Encode())

	cacheKey := cache.GenerateKey("recgov:facilities", map[string]string{
		"query":    params.Query,
		"state":    params.State,
		"activity": params.Activity,
		"limit":    strconv.Itoa(params.Limit),
		"offset":   strconv.Itoa(params.Offset),
	})

	if cached, found := c.cache.Get(cacheKey); found {
		c.logger.Debug("Cache hit for facilities search")
		var facilities []models.Facility
		if err := json.Unmarshal(cached, &facilities); err == nil {
			return facilities, nil
		}
	}

	headers := map[string]string{
		"apikey": c.apiKey,
	}

	resp, err := c.httpClient.Get(ctx, fullURL, headers)
	if err != nil {
		return nil, fmt.Errorf("failed to search facilities: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := util.ReadResponseBody(resp)
	if err != nil {
		return nil, err
	}

	var recgovResp models.RecreationGovResponse
	if err := json.Unmarshal(body, &recgovResp); err != nil {
		return nil, fmt.Errorf("failed to parse facilities response: %w", err)
	}

	facilitiesData, err := json.Marshal(recgovResp.RECDATA)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal facilities data: %w", err)
	}

	var facilities []models.Facility
	if err := json.Unmarshal(facilitiesData, &facilities); err != nil {
		return nil, fmt.Errorf("failed to unmarshal facilities: %w", err)
	}

	if cacheData, err := json.Marshal(facilities); err == nil {
		_ = c.cache.Set(cacheKey, cacheData)
	}

	c.logger.Info("Facilities search completed", "count", len(facilities))
	return facilities, nil
}

// GetFacilityDetails retrieves detailed information about a specific facility
func (c *recreationGovClient) GetFacilityDetails(ctx context.Context, facilityID string) (*models.Facility, error) {
	c.logger.Debug("Getting facility details", "facilityID", facilityID)

	endpoint := fmt.Sprintf("%s/facilities/%s", c.baseURL, facilityID)

	cacheKey := cache.GenerateKey("recgov:facility", map[string]string{"facilityID": facilityID})

	if cached, found := c.cache.Get(cacheKey); found {
		c.logger.Debug("Cache hit for facility details")
		var facility models.Facility
		if err := json.Unmarshal(cached, &facility); err == nil {
			return &facility, nil
		}
	}

	headers := map[string]string{
		"apikey": c.apiKey,
	}

	resp, err := c.httpClient.Get(ctx, endpoint, headers)
	if err != nil {
		return nil, fmt.Errorf("failed to get facility details: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := util.ReadResponseBody(resp)
	if err != nil {
		return nil, err
	}

	var facility models.Facility
	if err := json.Unmarshal(body, &facility); err != nil {
		return nil, fmt.Errorf("failed to unmarshal facility: %w", err)
	}

	if cacheData, err := json.Marshal(facility); err == nil {
		_ = c.cache.Set(cacheKey, cacheData)
	}

	c.logger.Info("Facility details retrieved", "facilityID", facilityID)
	return &facility, nil
}

// GetCampsites searches for campsites
func (c *recreationGovClient) GetCampsites(ctx context.Context, params SearchCampsitesParams) ([]models.Campsite, error) {
	c.logger.Debug("Searching campsites", "query", params.Query, "state", params.State)

	endpoint := fmt.Sprintf("%s/campsites", c.baseURL)
	queryParams := c.buildCampsitesQueryParams(params)
	fullURL := fmt.Sprintf("%s?%s", endpoint, queryParams.Encode())

	cacheKey := cache.GenerateKey("recgov:campsites", map[string]string{
		"query":  params.Query,
		"state":  params.State,
		"limit":  strconv.Itoa(params.Limit),
		"offset": strconv.Itoa(params.Offset),
	})

	if cached, found := c.cache.Get(cacheKey); found {
		c.logger.Debug("Cache hit for campsites search")
		var campsites []models.Campsite
		if err := json.Unmarshal(cached, &campsites); err == nil {
			return campsites, nil
		}
	}

	headers := map[string]string{
		"apikey": c.apiKey,
	}

	resp, err := c.httpClient.Get(ctx, fullURL, headers)
	if err != nil {
		return nil, fmt.Errorf("failed to search campsites: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := util.ReadResponseBody(resp)
	if err != nil {
		return nil, err
	}

	var recgovResp models.RecreationGovResponse
	if err := json.Unmarshal(body, &recgovResp); err != nil {
		return nil, fmt.Errorf("failed to parse campsites response: %w", err)
	}

	campsitesData, err := json.Marshal(recgovResp.RECDATA)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal campsites data: %w", err)
	}

	var campsites []models.Campsite
	if err := json.Unmarshal(campsitesData, &campsites); err != nil {
		return nil, fmt.Errorf("failed to unmarshal campsites: %w", err)
	}

	if cacheData, err := json.Marshal(campsites); err == nil {
		_ = c.cache.Set(cacheKey, cacheData)
	}

	c.logger.Info("Campsites search completed", "count", len(campsites))
	return campsites, nil
}

// GetActivities retrieves all available activities
func (c *recreationGovClient) GetActivities(ctx context.Context) ([]models.RecreationActivity, error) {
	c.logger.Debug("Getting activities")

	endpoint := fmt.Sprintf("%s/activities", c.baseURL)

	cacheKey := cache.GenerateKey("recgov:activities", map[string]string{})

	if cached, found := c.cache.Get(cacheKey); found {
		c.logger.Debug("Cache hit for activities")
		var activities []models.RecreationActivity
		if err := json.Unmarshal(cached, &activities); err == nil {
			return activities, nil
		}
	}

	headers := map[string]string{
		"apikey": c.apiKey,
	}

	resp, err := c.httpClient.Get(ctx, endpoint, headers)
	if err != nil {
		return nil, fmt.Errorf("failed to get activities: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := util.ReadResponseBody(resp)
	if err != nil {
		return nil, err
	}

	var recgovResp models.RecreationGovResponse
	if err := json.Unmarshal(body, &recgovResp); err != nil {
		return nil, fmt.Errorf("failed to parse activities response: %w", err)
	}

	activitiesData, err := json.Marshal(recgovResp.RECDATA)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal activities data: %w", err)
	}

	var activities []models.RecreationActivity
	if err := json.Unmarshal(activitiesData, &activities); err != nil {
		return nil, fmt.Errorf("failed to unmarshal activities: %w", err)
	}

	if cacheData, err := json.Marshal(activities); err == nil {
		_ = c.cache.Set(cacheKey, cacheData)
	}

	c.logger.Info("Activities retrieved", "count", len(activities))
	return activities, nil
}

// buildRecAreasQueryParams builds URL query parameters for recreation areas search
func (c *recreationGovClient) buildRecAreasQueryParams(params SearchRecAreasParams) url.Values {
	queryParams := url.Values{}

	if params.Query != "" {
		queryParams.Add("query", params.Query)
	}
	if params.State != "" {
		queryParams.Add("state", params.State)
	}
	if params.Activity != "" {
		queryParams.Add("activity", params.Activity)
	}
	if params.Limit > 0 {
		queryParams.Add("limit", strconv.Itoa(params.Limit))
	}
	if params.Offset > 0 {
		queryParams.Add("offset", strconv.Itoa(params.Offset))
	}

	return queryParams
}

// buildFacilitiesQueryParams builds URL query parameters for facilities search
func (c *recreationGovClient) buildFacilitiesQueryParams(params SearchFacilitiesParams) url.Values {
	queryParams := url.Values{}

	if params.Query != "" {
		queryParams.Add("query", params.Query)
	}
	if params.State != "" {
		queryParams.Add("state", params.State)
	}
	if params.Activity != "" {
		queryParams.Add("activity", params.Activity)
	}
	if params.Limit > 0 {
		queryParams.Add("limit", strconv.Itoa(params.Limit))
	}
	if params.Offset > 0 {
		queryParams.Add("offset", strconv.Itoa(params.Offset))
	}

	return queryParams
}

// buildCampsitesQueryParams builds URL query parameters for campsites search
func (c *recreationGovClient) buildCampsitesQueryParams(params SearchCampsitesParams) url.Values {
	queryParams := url.Values{}

	if params.Query != "" {
		queryParams.Add("query", params.Query)
	}
	if params.State != "" {
		queryParams.Add("state", params.State)
	}
	if params.Limit > 0 {
		queryParams.Add("limit", strconv.Itoa(params.Limit))
	}
	if params.Offset > 0 {
		queryParams.Add("offset", strconv.Itoa(params.Offset))
	}

	return queryParams
}
