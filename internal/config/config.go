// config/config.go
package config

import (
	"encoding/json"
	"os"
	"sync"
)

// Config holds all application configuration
type Config struct {
	AppName     string `json:"app_name"`
	Version     string `json:"version"`
	DefaultPage string `json:"default_page"`
	Theme       string `json:"theme"`
	Debug       bool   `json:"debug"`
	APIEndpoint string `json:"api_endpoint"`
	LogLevel    string `json:"log_level"`
}

var (
	config *Config
	once   sync.Once
)

// Load loads configuration from file or environment
func Load() *Config {
	once.Do(func() {
		config = &Config{
			AppName:     "tuix App",
			Version:     "1.0.0",
			DefaultPage: "home",
			Theme:       "dark",
			Debug:       false,
			APIEndpoint: "http://localhost:8080",
			LogLevel:    "info",
		}

		// Try to load from config file
		if data, err := os.ReadFile("config.json"); err == nil {
			var cfg Config
			if err := json.Unmarshal(data, &cfg); err == nil {
				config = &cfg
			}
		}

		// Override with environment variables
		if endpoint := os.Getenv("API_ENDPOINT"); endpoint != "" {
			config.APIEndpoint = endpoint
		}
		if debug := os.Getenv("DEBUG"); debug != "" {
			config.Debug = debug == "true"
		}
	})

	return config
}
