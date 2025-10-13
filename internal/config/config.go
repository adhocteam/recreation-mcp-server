package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"gopkg.in/yaml.v3"
)

// Config holds all application configuration
type Config struct {
	// API Keys
	NPSAPIKey           string
	RecreationGovAPIKey string
	OpenWeatherAPIKey   string

	// Server Configuration
	LogLevel             string
	CacheEnabled         bool
	CacheTTLSeconds      int
	MaxRequestsPerMinute int

	// API Configuration (optional, from YAML)
	APIs    APIConfig     `yaml:"apis"`
	Cache   CacheConfig   `yaml:"cache"`
	Logging LoggingConfig `yaml:"logging"`
}

// APIConfig holds configuration for external APIs
type APIConfig struct {
	NPS struct {
		BaseURL       string        `yaml:"base_url"`
		Timeout       time.Duration `yaml:"timeout"`
		RetryAttempts int           `yaml:"retry_attempts"`
	} `yaml:"nps"`
	RecreationGov struct {
		BaseURL       string        `yaml:"base_url"`
		Timeout       time.Duration `yaml:"timeout"`
		RetryAttempts int           `yaml:"retry_attempts"`
	} `yaml:"recreation_gov"`
	OpenWeather struct {
		BaseURL       string        `yaml:"base_url"`
		Timeout       time.Duration `yaml:"timeout"`
		RetryAttempts int           `yaml:"retry_attempts"`
	} `yaml:"openweather"`
}

// CacheConfig holds cache configuration
type CacheConfig struct {
	Enabled bool          `yaml:"enabled"`
	TTL     time.Duration `yaml:"ttl"`
	MaxSize string        `yaml:"max_size"`
}

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

// Load loads configuration from environment variables and optional YAML file
func Load() (*Config, error) {
	config := &Config{
		// Load required API keys
		NPSAPIKey:           os.Getenv("NPS_API_KEY"),
		RecreationGovAPIKey: os.Getenv("RECREATION_GOV_API_KEY"),
		OpenWeatherAPIKey:   os.Getenv("OPENWEATHER_API_KEY"),

		// Load optional configuration with defaults
		LogLevel:             getEnvWithDefault("LOG_LEVEL", "info"),
		CacheEnabled:         getEnvAsBool("CACHE_ENABLED", true),
		CacheTTLSeconds:      getEnvAsInt("CACHE_TTL_SECONDS", 3600),
		MaxRequestsPerMinute: getEnvAsInt("MAX_REQUESTS_PER_MINUTE", 60),
	}

	// Set default API configurations
	setDefaultAPIConfig(config)

	// Attempt to load YAML config if it exists
	if err := loadYAMLConfig(config, "config.yaml"); err != nil {
		// YAML is optional, only return error if file exists but can't be parsed
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("error loading config.yaml: %w", err)
		}
	}

	// Validate required configuration
	if err := config.Validate(); err != nil {
		return nil, err
	}

	return config, nil
}

// Validate checks that all required configuration is present
func (c *Config) Validate() error {
	if c.NPSAPIKey == "" {
		return fmt.Errorf("NPS_API_KEY is required")
	}
	if c.RecreationGovAPIKey == "" {
		return fmt.Errorf("RECREATION_GOV_API_KEY is required")
	}
	if c.OpenWeatherAPIKey == "" {
		return fmt.Errorf("OPENWEATHER_API_KEY is required")
	}

	// Validate log level
	validLogLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}
	if !validLogLevels[c.LogLevel] {
		return fmt.Errorf("invalid LOG_LEVEL: %s (must be debug, info, warn, or error)", c.LogLevel)
	}

	return nil
}

// setDefaultAPIConfig sets default values for API configurations
func setDefaultAPIConfig(config *Config) {
	// NPS defaults
	config.APIs.NPS.BaseURL = "https://developer.nps.gov/api/v1"
	config.APIs.NPS.Timeout = 30 * time.Second
	config.APIs.NPS.RetryAttempts = 3

	// Recreation.gov defaults
	config.APIs.RecreationGov.BaseURL = "https://ridb.recreation.gov/api/v1"
	config.APIs.RecreationGov.Timeout = 30 * time.Second
	config.APIs.RecreationGov.RetryAttempts = 3

	// OpenWeather defaults
	config.APIs.OpenWeather.BaseURL = "https://api.openweathermap.org/data/2.5"
	config.APIs.OpenWeather.Timeout = 15 * time.Second
	config.APIs.OpenWeather.RetryAttempts = 2

	// Cache defaults
	config.Cache.Enabled = config.CacheEnabled
	config.Cache.TTL = time.Duration(config.CacheTTLSeconds) * time.Second
	config.Cache.MaxSize = "100MB"

	// Logging defaults
	config.Logging.Level = config.LogLevel
	config.Logging.Format = "json"
}

// loadYAMLConfig loads configuration from a YAML file if it exists
func loadYAMLConfig(config *Config, filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, config)
}

// getEnvWithDefault gets an environment variable or returns a default value
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt gets an environment variable as an integer or returns a default value
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}

	return value
}

// getEnvAsBool gets an environment variable as a boolean or returns a default value
func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		return defaultValue
	}

	return value
}
