package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/adhocteam/recreation-mcp-server/internal/cache"
	"github.com/adhocteam/recreation-mcp-server/pkg/util"
)

// loadTestData loads a test fixture from testdata directory
func loadTestData(t *testing.T, filename string) []byte {
	t.Helper()
	data, err := os.ReadFile("../../testdata/" + filename)
	if err != nil {
		t.Fatalf("Failed to load test data %s: %v", filename, err)
	}
	return data
}

// setupTestNPSClient creates a test NPS client with mock server
func setupTestNPSClient(t *testing.T, handler http.HandlerFunc) (NPSClient, *httptest.Server) {
	t.Helper()
	server := httptest.NewServer(handler)

	logger := util.NewLogger("error", "json")                         // Use error level to reduce test noise
	httpClient := util.NewHTTPClient(100*time.Millisecond, 0, logger) // No retries for tests
	testCache := cache.NewCache(60*time.Second, 10, false)            // Disable cache for tests

	client := NewNPSClient(server.URL, "test-api-key", httpClient, testCache, logger)
	return client, server
}

func TestNPSClient_SearchParks(t *testing.T) {
	client, server := setupTestNPSClient(t, func(w http.ResponseWriter, r *http.Request) {
		// Verify API key is present
		if r.URL.Query().Get("api_key") == "" {
			t.Error("Expected api_key query parameter")
		}

		// Verify query parameters
		if stateCode := r.URL.Query().Get("stateCode"); stateCode != "CA" {
			t.Errorf("Expected stateCode=CA, got %s", stateCode)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(loadTestData(t, "nps_parks.json"))
	})
	defer server.Close()

	params := SearchParksParams{
		State: "CA",
		Limit: 10,
	}

	parks, err := client.SearchParks(context.Background(), params)
	if err != nil {
		t.Fatalf("SearchParks failed: %v", err)
	}

	if len(parks) != 2 {
		t.Errorf("Expected 2 parks, got %d", len(parks))
	}

	// Verify first park
	if parks[0].ParkCode != "yose" {
		t.Errorf("Expected park code 'yose', got '%s'", parks[0].ParkCode)
	}
	if parks[0].Name != "Yosemite" {
		t.Errorf("Expected park name 'Yosemite', got '%s'", parks[0].Name)
	}
	if parks[0].States != "CA" {
		t.Errorf("Expected state 'CA', got '%s'", parks[0].States)
	}
}

func TestNPSClient_GetParkDetails(t *testing.T) {
	client, server := setupTestNPSClient(t, func(w http.ResponseWriter, r *http.Request) {
		// Verify endpoint
		if r.URL.Path != "/parks" {
			t.Errorf("Expected path /parks, got %s", r.URL.Path)
		}

		// Verify parkCode parameter
		if parkCode := r.URL.Query().Get("parkCode"); parkCode != "yose" {
			t.Errorf("Expected parkCode=yose, got %s", parkCode)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(loadTestData(t, "nps_parks.json"))
	})
	defer server.Close()

	park, err := client.GetParkDetails(context.Background(), "yose")
	if err != nil {
		t.Fatalf("GetParkDetails failed: %v", err)
	}

	if park.ParkCode != "yose" {
		t.Errorf("Expected park code 'yose', got '%s'", park.ParkCode)
	}
	if park.FullName != "Yosemite National Park" {
		t.Errorf("Expected full name 'Yosemite National Park', got '%s'", park.FullName)
	}

	// Verify nested data
	if len(park.Activities) != 2 {
		t.Errorf("Expected 2 activities, got %d", len(park.Activities))
	}
	if len(park.EntranceFees) != 1 {
		t.Errorf("Expected 1 entrance fee, got %d", len(park.EntranceFees))
	}
	if len(park.Hours) != 1 {
		t.Errorf("Expected 1 operating hours entry, got %d", len(park.Hours))
	}
}

func TestNPSClient_GetAlerts(t *testing.T) {
	client, server := setupTestNPSClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/alerts" {
			t.Errorf("Expected path /alerts, got %s", r.URL.Path)
		}

		if parkCode := r.URL.Query().Get("parkCode"); parkCode != "yose" {
			t.Errorf("Expected parkCode=yose, got %s", parkCode)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(loadTestData(t, "nps_alerts.json"))
	})
	defer server.Close()

	alerts, err := client.GetAlerts(context.Background(), "yose")
	if err != nil {
		t.Fatalf("GetAlerts failed: %v", err)
	}

	if len(alerts) != 2 {
		t.Errorf("Expected 2 alerts, got %d", len(alerts))
	}

	// Verify first alert
	if alerts[0].Category != "Park Closure" {
		t.Errorf("Expected category 'Park Closure', got '%s'", alerts[0].Category)
	}
	if alerts[0].Title != "Tioga Road Closed for Winter" {
		t.Errorf("Expected title 'Tioga Road Closed for Winter', got '%s'", alerts[0].Title)
	}
}

func TestNPSClient_GetCampgrounds(t *testing.T) {
	client, server := setupTestNPSClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/campgrounds" {
			t.Errorf("Expected path /campgrounds, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(loadTestData(t, "nps_campgrounds.json"))
	})
	defer server.Close()

	params := SearchCampgroundsParams{
		Query: "Upper Pines",
		Limit: 10,
	}

	campgrounds, err := client.GetCampgrounds(context.Background(), params)
	if err != nil {
		t.Fatalf("GetCampgrounds failed: %v", err)
	}

	if len(campgrounds) != 1 {
		t.Errorf("Expected 1 campground, got %d", len(campgrounds))
	}

	// Verify campground data
	if campgrounds[0].Name != "Upper Pines Campground" {
		t.Errorf("Expected name 'Upper Pines Campground', got '%s'", campgrounds[0].Name)
	}
	if campgrounds[0].ParkCode != "yose" {
		t.Errorf("Expected park code 'yose', got '%s'", campgrounds[0].ParkCode)
	}
}

func TestNPSClient_GetActivities(t *testing.T) {
	client, server := setupTestNPSClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/activities" {
			t.Errorf("Expected path /activities, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(loadTestData(t, "nps_activities.json"))
	})
	defer server.Close()

	activities, err := client.GetActivities(context.Background())
	if err != nil {
		t.Fatalf("GetActivities failed: %v", err)
	}

	if len(activities) != 5 {
		t.Errorf("Expected 5 activities, got %d", len(activities))
	}

	// Verify activities
	expectedActivities := []string{"Astronomy", "Biking", "Boating", "Camping", "Hiking"}
	for i, activity := range activities {
		if activity.Name != expectedActivities[i] {
			t.Errorf("Expected activity '%s', got '%s'", expectedActivities[i], activity.Name)
		}
	}
}

func TestNPSClient_ErrorHandling(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		wantError  bool
		wantEmpty  bool
	}{
		{"Success", http.StatusOK, false, false},
		{"Bad Request", http.StatusBadRequest, false, true}, // 4xx returns empty results, not error
		{"Unauthorized", http.StatusUnauthorized, false, true},
		{"Not Found", http.StatusNotFound, false, true},
		{"Server Error", http.StatusInternalServerError, true, false}, // 5xx returns error
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, server := setupTestNPSClient(t, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				if tt.statusCode == http.StatusOK {
					w.Write(loadTestData(t, "nps_parks.json"))
				} else {
					// For 4xx, return empty NPS response; for 5xx, return invalid JSON
					if tt.statusCode < 500 {
						w.Write([]byte(`{"total": "0", "data": []}`))
					} else {
						w.Write([]byte(`{"error": "Server error"}`))
					}
				}
			})
			defer server.Close()

			parks, err := client.SearchParks(context.Background(), SearchParksParams{})

			if tt.wantError && err == nil {
				t.Error("Expected error, got nil")
			}
			if !tt.wantError && err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
			if tt.wantEmpty && len(parks) != 0 {
				t.Errorf("Expected empty results, got %d parks", len(parks))
			}
		})
	}
}

func TestNPSClient_ContextCancellation(t *testing.T) {
	client, server := setupTestNPSClient(t, func(w http.ResponseWriter, r *http.Request) {
		// Simulate slow response - check if context is cancelled
		select {
		case <-r.Context().Done():
			return
		case <-time.After(100 * time.Millisecond):
			w.Write(loadTestData(t, "nps_parks.json"))
		}
	})
	defer server.Close()

	// Create a context that's already cancelled
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := client.SearchParks(ctx, SearchParksParams{})
	if err == nil {
		t.Error("Expected context cancellation error, got nil")
	}
}
