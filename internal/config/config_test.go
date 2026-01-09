package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	// Set test environment variables
	os.Setenv("NPS_API_KEY", "test-nps-key")
	os.Setenv("RECREATION_GOV_API_KEY", "test-recgov-key")
	os.Setenv("OPENWEATHER_API_KEY", "test-weather-key")
	os.Setenv("LOG_LEVEL", "debug")
	os.Setenv("CACHE_ENABLED", "true")
	os.Setenv("CACHE_TTL_SECONDS", "7200")
	os.Setenv("MAX_REQUESTS_PER_MINUTE", "120")
	defer func() {
		os.Unsetenv("NPS_API_KEY")
		os.Unsetenv("RECREATION_GOV_API_KEY")
		os.Unsetenv("OPENWEATHER_API_KEY")
		os.Unsetenv("LOG_LEVEL")
		os.Unsetenv("CACHE_ENABLED")
		os.Unsetenv("CACHE_TTL_SECONDS")
		os.Unsetenv("MAX_REQUESTS_PER_MINUTE")
	}()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Verify API keys
	if cfg.NPSAPIKey != "test-nps-key" {
		t.Errorf("Expected NPS key 'test-nps-key', got '%s'", cfg.NPSAPIKey)
	}
	if cfg.RecreationGovAPIKey != "test-recgov-key" {
		t.Errorf("Expected RecGov key 'test-recgov-key', got '%s'", cfg.RecreationGovAPIKey)
	}
	if cfg.OpenWeatherAPIKey != "test-weather-key" {
		t.Errorf("Expected Weather key 'test-weather-key', got '%s'", cfg.OpenWeatherAPIKey)
	}

	// Verify other settings
	if cfg.LogLevel != "debug" {
		t.Errorf("Expected log level 'debug', got '%s'", cfg.LogLevel)
	}
	if !cfg.CacheEnabled {
		t.Error("Expected cache to be enabled")
	}
	if cfg.CacheTTLSeconds != 7200 {
		t.Errorf("Expected cache TTL 7200, got %d", cfg.CacheTTLSeconds)
	}
	if cfg.MaxRequestsPerMinute != 120 {
		t.Errorf("Expected max requests 120, got %d", cfg.MaxRequestsPerMinute)
	}
}

func TestLoad_MissingAPIKeys(t *testing.T) {
	// Clear all API keys
	os.Unsetenv("NPS_API_KEY")
	os.Unsetenv("RECREATION_GOV_API_KEY")
	os.Unsetenv("OPENWEATHER_API_KEY")

	_, err := Load()
	if err == nil {
		t.Error("Expected error when API keys are missing")
	}
}

func TestLoad_Defaults(t *testing.T) {
	// Set only required API keys
	os.Setenv("NPS_API_KEY", "test-key-1")
	os.Setenv("RECREATION_GOV_API_KEY", "test-key-2")
	os.Setenv("OPENWEATHER_API_KEY", "test-key-3")
	// Clear optional settings
	os.Unsetenv("LOG_LEVEL")
	os.Unsetenv("CACHE_ENABLED")
	os.Unsetenv("CACHE_TTL_SECONDS")
	os.Unsetenv("MAX_REQUESTS_PER_MINUTE")
	defer func() {
		os.Unsetenv("NPS_API_KEY")
		os.Unsetenv("RECREATION_GOV_API_KEY")
		os.Unsetenv("OPENWEATHER_API_KEY")
	}()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Verify defaults
	if cfg.LogLevel != "info" {
		t.Errorf("Expected default log level 'info', got '%s'", cfg.LogLevel)
	}
	if !cfg.CacheEnabled {
		t.Error("Expected cache to be enabled by default")
	}
	if cfg.CacheTTLSeconds != 3600 {
		t.Errorf("Expected default cache TTL 3600, got %d", cfg.CacheTTLSeconds)
	}
	if cfg.MaxRequestsPerMinute != 60 {
		t.Errorf("Expected default max requests 60, got %d", cfg.MaxRequestsPerMinute)
	}
}

func TestLoad_InvalidBooleans(t *testing.T) {
	os.Setenv("NPS_API_KEY", "test-key")
	os.Setenv("RECREATION_GOV_API_KEY", "test-key")
	os.Setenv("OPENWEATHER_API_KEY", "test-key")
	os.Setenv("CACHE_ENABLED", "invalid")
	defer func() {
		os.Unsetenv("NPS_API_KEY")
		os.Unsetenv("RECREATION_GOV_API_KEY")
		os.Unsetenv("OPENWEATHER_API_KEY")
		os.Unsetenv("CACHE_ENABLED")
	}()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load should not fail on invalid boolean: %v", err)
	}

	// Should use default value
	if !cfg.CacheEnabled {
		t.Error("Expected cache to use default value (true) when invalid boolean provided")
	}
}

func TestLoad_InvalidNumbers(t *testing.T) {
	os.Setenv("NPS_API_KEY", "test-key")
	os.Setenv("RECREATION_GOV_API_KEY", "test-key")
	os.Setenv("OPENWEATHER_API_KEY", "test-key")
	os.Setenv("CACHE_TTL_SECONDS", "not-a-number")
	os.Setenv("MAX_REQUESTS_PER_MINUTE", "also-not-a-number")
	defer func() {
		os.Unsetenv("NPS_API_KEY")
		os.Unsetenv("RECREATION_GOV_API_KEY")
		os.Unsetenv("OPENWEATHER_API_KEY")
		os.Unsetenv("CACHE_TTL_SECONDS")
		os.Unsetenv("MAX_REQUESTS_PER_MINUTE")
	}()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load should not fail on invalid numbers: %v", err)
	}

	// Should use default values
	if cfg.CacheTTLSeconds != 3600 {
		t.Errorf("Expected default TTL 3600, got %d", cfg.CacheTTLSeconds)
	}
	if cfg.MaxRequestsPerMinute != 60 {
		t.Errorf("Expected default max requests 60, got %d", cfg.MaxRequestsPerMinute)
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name      string
		config    Config
		wantError bool
	}{
		{
			name: "Valid config",
			config: Config{
				NPSAPIKey:            "key1",
				RecreationGovAPIKey:  "key2",
				OpenWeatherAPIKey:    "key3",
				LogLevel:             "info",
				CacheEnabled:         true,
				CacheTTLSeconds:      3600,
				MaxRequestsPerMinute: 60,
			},
			wantError: false,
		},
		{
			name: "Missing NPS key",
			config: Config{
				RecreationGovAPIKey: "key2",
				OpenWeatherAPIKey:   "key3",
				LogLevel:            "info",
			},
			wantError: true,
		},
		{
			name: "Missing RecGov key",
			config: Config{
				NPSAPIKey:         "key1",
				OpenWeatherAPIKey: "key3",
				LogLevel:          "info",
			},
			wantError: true,
		},
		{
			name: "Missing Weather key",
			config: Config{
				NPSAPIKey:           "key1",
				RecreationGovAPIKey: "key2",
				LogLevel:            "info",
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.wantError && err == nil {
				t.Error("Expected validation error, got nil")
			}
			if !tt.wantError && err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
		})
	}
}
