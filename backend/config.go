package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config holds the application configuration loaded from config.json.
type Config struct {
	Port string `json:"port"`
}

// LoadConfig reads and parses the JSON configuration file at the given path.
func LoadConfig(path string) (*Config, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", path, err)
	}

	var cfg Config
	if err := json.Unmarshal(file, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Default port if not set
	if cfg.Port == "" {
		cfg.Port = "8080"
	}

	return &cfg, nil
}
