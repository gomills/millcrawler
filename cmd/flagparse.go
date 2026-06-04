package main

import (
	"errors"
	"flag"
	"strings"

	"github.com/gomills/millcrawler/pkg/config"
)

// parseFlagsIntoConfig parses flags from the CLI and instantiates a new config for the crawler. A friendly user is expected to call this function.
// In other words, it doesn't have unbreakable validation, just basic.
func parseFlagsIntoConfig() (string, *config.Config, error) {

	// load config from flags
	domain := flag.String("domain", "", "Target domain to crawl. Example: 'example.com'")
	timeOutDuration := flag.Int("timeoutseconds", 60, "Maximum crawl duration in seconds. Example: '60'")
	maxPathDepth := flag.Int("maxpathdepth", 1, "Maximum URL path depth to crawl. Example: '2' allows '/path1/path2'")
	maxNCrawledUrls := flag.Int("maxnumurls", 100, "Maximum number of URLs to crawl before stopping")
	allowedExternalDomains := flag.String("allowedextdomains", "", "Comma-separated list of allowed external domains. Example: 'github.com,gitlab.com'")
	sensitivePatterns := flag.String("sensitivepatterns", "", "Comma-separated URL patterns that bypass max path depth restrictions. Example: 'api,admin,debug'")
	allowedExtensions := flag.String("allowedextensions", ",.htm,.html,.js", "Comma-separated list of allowed file extensions. Example: '.html,.js,.json'")
	workers := flag.Int("workers", 1, "Number of concurrent crawling workers")
	printUrls := flag.Bool("printurls", false, "Interactive logging: banner, URLs, to stdout.")
	bruteforceSubdomains := flag.Bool("bruteforcesubdomains", true, "Enable subdomain bruteforcing. E.g: 'true' for testing.example.com, admin.example.com, ...")
	scanSecrets := flag.Bool("scansecrets", false, "Decide if scan responses for secrets")
	cookies := flag.Bool("cookies", false, "Enable cookie jar for use during crawling")
	flag.Parse()

	// manually validate the only value not used for config where there's native validation
	if *domain == "" {
		return "", nil, errors.New("Domain is empty")
	}

	configuration, err := config.LoadConfig(
		*timeOutDuration,
		*maxPathDepth,
		*maxNCrawledUrls,
		strings.Split(*allowedExternalDomains, ","),
		strings.Split(*sensitivePatterns, ","),
		strings.Split(*allowedExtensions, ","),
		*workers,
		*printUrls,
		*bruteforceSubdomains,
		*scanSecrets,
		*cookies,
	)

	return *domain, configuration, err
}
