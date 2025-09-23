package urlvalidator

import (
	"fmt"
	"net/url"
	"slices"

	"github.com/gomills/gofocusedcrawler/internal/config"
)

// validateExternalUrl checks if external url's domain is in the allowed domains from config
func validateExternalUrl(config *config.Config, parsedUrl *url.URL) (*url.URL, error) {
	if slices.Contains(config.AllowedExternalDomains, parsedUrl.Hostname()) {
		return parsedUrl, nil
	}
	return nil, fmt.Errorf("External URL not in allowed domains")
}
