package urlvalidator

import (
	"fmt"
	"net/url"
	"slices"

	"github.com/gomills/millcrawler/pkg/config"
)

// validateExternalUrl checks if external url's domain is in the allowed domains from config
func validateExternalUrl(config *config.Config, parsedUrl *url.URL) (*url.URL, error) {
	if slices.Contains(config.AllowedExternalDomains, parsedUrl.Hostname()) {
		return parsedUrl, nil
	}
	return nil, fmt.Errorf("external url doesn't come from any allowed external domain")
}
