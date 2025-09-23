package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	yaml "sigs.k8s.io/yaml/goyaml.v2"
)

// RawConfig holds the configuration as it appears in the file, with comma-separated strings
type RawConfig struct {
	TimeOutDuration        int    `yaml:"timeout_duration_sec" json:"timeout_duration_sec"`
	CustomUserAgent        string `yaml:"custom_user_agent" json:"custom_user_agent"`
	AllowedExternalDomains string `yaml:"allowed_external_domains" json:"allowed_external_domains"`
	MaxPathDepth           int    `yaml:"max_path_depth" json:"max_path_depth"`
	SensitivePatterns      string `yaml:"sensitive_patterns" json:"sensitive_patterns"`
	AllowedExtensions      string `yaml:"allowed_extensions" json:"allowed_extensions"`
	Workers                int    `yaml:"workers" json:"workers"`
}

// Config holds the processed configuration with proper slices
type Config struct {
	TimeOutDuration        int      `yaml:"timeout_duration_sec" json:"timeout_duration_sec"`
	CustomUserAgent        string   `yaml:"custom_user_agent" json:"custom_user_agent"`
	AllowedExternalDomains []string `yaml:"allowed_external_domains" json:"allowed_external_domains"`
	MaxPathDepth           int      `yaml:"max_path_depth" json:"max_path_depth"`
	SensitivePatterns      []string `yaml:"sensitive_patterns" json:"sensitive_patterns"`
	AllowedExtensions      []string `yaml:"allowed_extensions" json:"allowed_extensions"`
	Workers                int      `yaml:"workers" json:"workers"`
}

func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var rawConfig RawConfig

	if strings.HasSuffix(filename, ".json") {
		err = json.Unmarshal(data, &rawConfig)
	} else if strings.HasSuffix(filename, ".yaml") || strings.HasSuffix(filename, ".yml") {
		err = yaml.Unmarshal(data, &rawConfig)
	} else {
		return nil, fmt.Errorf("unsupported config format")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// Convert raw config to processed config
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
func parseCommaSeparated(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		result = append(result, part)
	}
	return result
}
