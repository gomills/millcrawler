package utils

import (
	"fmt"
	"net/url"
	"strings"
)

// GetBruteForcedUrls generates a list of brute-forced URLs to crawl from the provided domain.
// It creates URLs for common endpoints and subdomains that might be of interest.
// Returns a slice of *url.URL and an error if any url parsing fails.
func GetBruteForcedUrls(domain string, registeredDomain string) ([]*url.URL, error) {

	if strings.TrimSpace(domain) == "" || strings.TrimSpace(registeredDomain) == "" {
		return nil, fmt.Errorf("empty_domain_or_registeredDomain")
	}

	// build brute-forced urls as strings
	bruteForcedUrlStr := []string{

		// from domain
		fmt.Sprintf("https://%s/", domain),
		fmt.Sprintf("https://%s/robots.txt", domain),
		fmt.Sprintf("https://%s/sitemap.xml", domain),

		// from registered domain, for specific subdomains
		fmt.Sprintf("https://dev.%s/", registeredDomain),
		fmt.Sprintf("https://staging.%s/", registeredDomain),
		fmt.Sprintf("https://admin.%s/", registeredDomain),
		fmt.Sprintf("https://test.%s/", registeredDomain),
		fmt.Sprintf("https://internal.%s/", registeredDomain),
	}

	// parse brute-forced url strings and store them
	bruteForcedUrls := make([]*url.URL, len(bruteForcedUrlStr))

	for ind := range bruteForcedUrlStr {

		parsedUrl, err := url.Parse(bruteForcedUrlStr[ind])
		if err != nil {
			return nil, fmt.Errorf("seed_gen_failed")
		}

		bruteForcedUrls[ind] = parsedUrl

	}

	return bruteForcedUrls, nil

}
