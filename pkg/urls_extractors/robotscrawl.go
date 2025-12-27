package urls_extractors

import (
	"bufio"
	"context"
	"strings"

	"github.com/gomills/gofocusedcrawler/internal/config"
	"github.com/gomills/gofocusedcrawler/pkg/queue"
	"github.com/gomills/gofocusedcrawler/pkg/url_validator"
)

// extractUrlsFromRobots performs a URL extraction from robots by parsing it. This functions extracts ALL URLs, it doesn't respect ROBOTS.txt at all.
// In fact, it does the opposite
func extractUrlsFromRobots(ctx context.Context, bodyStr string, config *config.Config, qp *queue.Queue, domain string, registeredDomain string) {
	scanner := bufio.NewScanner(strings.NewReader(bodyStr))

	for scanner.Scan() {
		line := scanner.Text()
		cleanLine := strings.TrimSpace(line)

		if cleanLine == "" || strings.HasPrefix(cleanLine, "#") {
			continue
		}

		cleanLine = strings.ToLower(cleanLine)
		if strings.HasPrefix(cleanLine, "disallow:") || strings.HasPrefix(cleanLine, "allow:") || strings.HasPrefix(cleanLine, "sitemap:") {
			values := strings.SplitN(cleanLine, ":", 2)
			if len(values) == 1 {
				continue
			}

			values[1] = strings.TrimSpace(values[1])
			parsedUrl, err := url_validator.ValidateStringForUrl(config, values[1], domain, registeredDomain)
			if err == nil {
				qp.AddUrl(ctx, parsedUrl)
			}
		}
	}
}
