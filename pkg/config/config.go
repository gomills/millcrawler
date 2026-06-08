package config

import "fmt"

// Config holds the crawler configuration: parameters, heuristics, etc
type Config struct {
	TimeOutDuration        int // has no sentinel value, it must be specified always
	AllowedExternalDomains []string
	MaxPathDepth           int // sentinel value of 0 for unlimited
	MaxNCrawledUrls        int // sentinel value of 0 for unlimited
	SensitivePatterns      []string
	AllowedExtensions      []string
	Workers                int
	PrintUrls              bool
	BruteforceSubdomains   bool
	ScanSecrets            bool
	Cookies                bool
}

// LoadConfig returns the crawler config
func LoadConfig(
	timeoutSeconds int,
	maxPathDepth int,
	maxNCrawledUrls int,
	allowedExtDom []string,
	sensitivePatterns []string,
	allowedExtensions []string,
	workers int,
	printUrls bool,
	bruteforceSubdomains bool,
	scanSecrets bool,
	cookies bool) (*Config, error) {

	if timeoutSeconds < 1 {
		return nil, fmt.Errorf("timeoutSeconds must be 1 or greater, got %d", timeoutSeconds)
	}
	if maxPathDepth < 0 {
		return nil, fmt.Errorf("maxPathDepth must be 0 or greater, got %d", maxPathDepth)
	}
	if maxNCrawledUrls < 0 {
		return nil, fmt.Errorf("maxNCrawledUrls must be 0 or greater, got %d", maxNCrawledUrls)
	}
	if workers < 1 {
		return nil, fmt.Errorf("workers must be 1 or greater, got %d", workers)
	}

	return &Config{
		TimeOutDuration:        timeoutSeconds,
		MaxPathDepth:           maxPathDepth,
		MaxNCrawledUrls:        maxNCrawledUrls,
		AllowedExternalDomains: allowedExtDom,
		SensitivePatterns:      sensitivePatterns,
		AllowedExtensions:      allowedExtensions,
		Workers:                workers,
		PrintUrls:              printUrls,
		BruteforceSubdomains:   bruteforceSubdomains,
		ScanSecrets:            scanSecrets,
		Cookies:                cookies,
	}, nil
}
