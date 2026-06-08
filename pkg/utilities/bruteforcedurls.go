package utilities

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/gomills/millcrawler/pkg/config"
)

// GetBruteForcedUrls generates a list of brute-forced URLs. It uses domain for bruteforcing
// paths and registered domain for bruteforcing subdomains.
// URLs are for common endpoints and subdomains that might be of interest.
// Returns a slice of *url.URL and an error if and if any url parsing fails.
func GetBruteForcedUrls(domain string, registeredDomain string, config *config.Config) ([]*url.URL, error) {

	if strings.TrimSpace(domain) == "" || strings.TrimSpace(registeredDomain) == "" {
		return nil, fmt.Errorf("empty domain or registered domain")
	}

	// build brute-forced urls as strings
	bruteForcedUrlStr := []string{

		// bruteforced paths
		fmt.Sprintf("https://%s/", domain),
		fmt.Sprintf("https://%s/robots.txt", domain),
		fmt.Sprintf("https://%s/sitemap.xml", domain),
	}

	// add subdomains if specified from config
	if config.BruteforceSubdomains {
		bruteForcedSubdomainsStr := []string{
			fmt.Sprintf("https://dev.%s/", registeredDomain),
			fmt.Sprintf("https://staging.%s/", registeredDomain),
			fmt.Sprintf("https://admin.%s/", registeredDomain),
			fmt.Sprintf("https://test.%s/", registeredDomain),
			fmt.Sprintf("https://internal.%s/", registeredDomain),
		}

		bruteForcedUrlStr = append(bruteForcedUrlStr, bruteForcedSubdomainsStr...)
	}

	// parse brute-forced url strings and store them
	bruteForcedUrls := make([]*url.URL, len(bruteForcedUrlStr))

	for ind := range bruteForcedUrlStr {

		parsedUrl, err := url.Parse(bruteForcedUrlStr[ind])
		if err != nil {
			return nil, fmt.Errorf("parsing of bruteforced url '%s' failed", bruteForcedUrlStr[ind])
		}

		bruteForcedUrls[ind] = parsedUrl

	}

	return bruteForcedUrls, nil

}
