package main

import (
	"errors"
	"flag"
	"strings"

	"github.com/gomills/millcrawler/pkg/config"
)

const (
	defaultAllowedExtDomains = "github.com,gitlab.com,docker.com,pastebin.com,nuget.org,bitbucket.org,s3.amazonaws.com,telegram.org,slack.com,drive.google.com,docs.google.com,codeberg.org"
	defaultAllowedExtensions = ",.js,.ts,.map,.properties,.log,.cfg,.pem,.crt,.npmrc,.yarnrc,.html,.htm,.json,.yaml,.yml,.env,.conf,.xml,.txt,.py,.php,.git,.sh,.key,.go,.ini,.example,.md,.rb,.java,.cpp,.c,.pl,.zsh,.bak,.old,.sql,.db,.tar,.gz,.zip,.ovpn,.toml"
)

// parseFlagsIntoConfig parses flags from the CLI and instantiates a new config for the crawler.
// It supposes that a friendly user will call this function (it doesn't have unbreakable validation, just basic).
func parseFlagsIntoConfig() (string, *config.Config, error) {

	// load config from flags
	domain := flag.String("domain", "", "Target domain to crawl. Example: -domain=example.com")
	timeOutDuration := flag.Int("timeoutseconds", 3600, "Maximum crawl duration in seconds, never unlimited. Example: -timeoutseconds=60")
	maxPathDepth := flag.Int("maxpathdepth", 0, "Maximum URL path depth to crawl, 0 for unlimited. Example: -maxpathdepth=2 allows /path1/path2 but not /path1/path2/path3")
	maxNCrawledUrls := flag.Int("maxnumurls", 0, "Maximum number of URLs to crawl before stopping, 0 for unlimited. Example: -maxnumurls=400")
	allowedExternalDomains := flag.String("allowedextdomains", defaultAllowedExtDomains, "Comma-separated list of allowed external domains. Example: -allowedextdomains=github.com,gitlab.com")
	sensitivePatterns := flag.String("sensitivepatterns", "", "Comma-separated URL patterns that bypass max path depth restrictions. Example: sensitivepatterns=api,admin,debug")
	allowedExtensions := flag.String("allowedextensions", defaultAllowedExtensions, "Comma-separated list of allowed file extensions ('' for subpages...). Example: allowedextensions=,.html,.js,.json")
	workers := flag.Int("workers", 1, "Number of concurrent crawling workers. Example: -workers=10")
	printUrls := flag.Bool("printurls", false, "Interactive logging: banner and URLs to stdout. If disabled, it outputs .json results directly for scripting integration. Example: -printurls=true")
	bruteforceSubdomains := flag.Bool("bruteforcesubdomains", false, "Enable subdomain bruteforcing. Example: -bruteforcesubdomains=true")
	scanSecrets := flag.Bool("scansecrets", false, "Decide if scan responses for secrets. Example: -scansecrets=true")
	cookies := flag.Bool("cookies", false, "Enable cookie jar for use during crawling. Example: -cookies=true")
	flag.Parse()

	// manually validate the only value not used for config (in config there's native validation already...)
	if *domain == "" {
		return "", nil, errors.New("Domain is empty")
	}

	// instantiate config and check for error where there's validation for the rest of the input parameters
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
