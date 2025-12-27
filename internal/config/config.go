package config

// Config holds the crawler configuration: parameters, heuristics, etc
type Config struct {
	TimeOutDuration        int
	AllowedExternalDomains []string
	MaxPathDepth           int
	SensitivePatterns      []string
	AllowedExtensions      []string
	Workers                int
}

// LoadConfig returns the crawler config
func LoadConfig() *Config {

	return &Config{
		TimeOutDuration:        90,
		MaxPathDepth:           3,
		AllowedExternalDomains: []string{"github.com", "bitbucket.com"},
		SensitivePatterns:      []string{"internal", "secret", "private"},
		AllowedExtensions:      []string{".js", ".html", ".htm", ".txt", ".py", ".php", ".git", ".json", ".yaml"},
		Workers:                5,
	}

}
