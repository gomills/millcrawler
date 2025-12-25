package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// RawConfig holds the configuration as it appears in the file, with comma-separated strings
type RawConfig struct {
	TimeOutDuration        int    `json:"timeout_duration_sec"`
	CustomUserAgent        string `json:"custom_user_agent"`
	AllowedExternalDomains string `json:"allowed_external_domains"`
	MaxPathDepth           int    `json:"max_path_depth"`
	SensitivePatterns      string `json:"sensitive_patterns"`
	AllowedExtensions      string `json:"allowed_extensions"`
	Workers                int    `json:"workers"`
}

// Config holds the processed configuration with proper go objects
type Config struct {
	TimeOutDuration        int      `json:"timeout_duration_sec"`
	CustomUserAgent        string   `json:"custom_user_agent"`
	AllowedExternalDomains []string `json:"allowed_external_domains"`
	MaxPathDepth           int      `json:"max_path_depth"`
	SensitivePatterns      []string `json:"sensitive_patterns"`
	AllowedExtensions      []string `json:"allowed_extensions"`
	Workers                int      `json:"workers"`
}

// LoadConfig reads the config.json and unmarshals it into a go object
func LoadConfig(filename string) (*Config, error) {

	// 1. read config.json
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// 2. unmarshal it into raw config
	var rawConfig RawConfig

	err = json.Unmarshal(data, &rawConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// 3. convert raw config to processed config with go objects
	config := &Config{
		TimeOutDuration:        rawConfig.TimeOutDuration,
		CustomUserAgent:        rawConfig.CustomUserAgent,
		MaxPathDepth:           rawConfig.MaxPathDepth,
		AllowedExternalDomains: parseCommaSeparated(rawConfig.AllowedExternalDomains),
		SensitivePatterns:      parseCommaSeparated(rawConfig.SensitivePatterns),
		AllowedExtensions:      parseCommaSeparated(rawConfig.AllowedExtensions),
		Workers:                rawConfig.Workers,
	}

	return config, nil
}

// parseCommaSeparated splits a comma-separated string into a slice of strings,
// trimming whitespace from each element
func parseCommaSeparated(str string) []string {
	if str == "" {
		return nil
	}

	// split string
	parts := strings.Split(str, ",")

	// populate corresponding slice with whitespace trimmed elements
	result := make([]string, len(parts))

	for ind := range parts {
		result[ind] = strings.TrimSpace(parts[ind])
	}

	return result
}
