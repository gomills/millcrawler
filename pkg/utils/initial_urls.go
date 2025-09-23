package utils

import (
	"fmt"
	"net/url"
	"strings"
)

// GetInitialUrls generates a list of brute-forced URLs to crawl based on the provided domain and registered domain.
// It creates URLs for common endpoints and subdomains that might be of interest.
// Returns a slice of parsed Urls and an error if any URL parsing fails.
func GetInitialUrls(domain string, registeredDomain string) ([]*url.URL, error) {

	if strings.TrimSpace(domain) == "" || strings.TrimSpace(registeredDomain) == "" {
		return nil, fmt.Errorf("domain or registeredDomain cannot be empty")
	}

	urlList := []string{
		fmt.Sprintf("https://%s/", domain),
		fmt.Sprintf("https://%s/robots.txt", domain),
		fmt.Sprintf("https://%s/sitemap.xml", domain),
		fmt.Sprintf("https://dev.%s/", registeredDomain),
		fmt.Sprintf("https://staging.%s/", registeredDomain),
		fmt.Sprintf("https://admin.%s/", registeredDomain),
		fmt.Sprintf("https://test.%s/", registeredDomain),
		fmt.Sprintf("https://internal.%s/", registeredDomain),
	}

	var toCrawlUrls []*url.URL

	for _, x := range urlList {

		parsedUrl, err := url.Parse(x)
		if err != nil {
			return nil, fmt.Errorf("failed to parse URL %q: %w", x, err)
		}

		toCrawlUrls = append(toCrawlUrls, parsedUrl)
	}

	return toCrawlUrls, nil

}
