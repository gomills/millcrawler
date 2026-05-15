package crawler

type CrawlingOutcome struct {
	// metrics
	Domain          string
	NumURLs         int
	DurationSeconds float64
	StopReason      string
	SecretsFound    int

	// secrets
	SecretsMap map[string]string
}

// GetCrawlingOutcome returns the crawling outcome, which consists in metrics and secrets from the crawling.
func CreateCrawlingOutcome(domain string, nURLs int, durationSeconds float64, stopReason string, secretsFound int, secretsMap map[string]string) *CrawlingOutcome {
	return &CrawlingOutcome{
		Domain:          domain,
		NumURLs:         nURLs,
		DurationSeconds: durationSeconds,
		StopReason:      stopReason,
		SecretsFound:    secretsFound,
		SecretsMap:      secretsMap,
	}
}
