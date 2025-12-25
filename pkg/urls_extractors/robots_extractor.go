package urls_extractors

import (
	"bufio"
	"context"
	"strings"

	"github.com/gomills/gofocusedcrawler/internal/config"
	"github.com/gomills/gofocusedcrawler/pkg/queue"
	"github.com/gomills/gofocusedcrawler/pkg/urlvalidator"
)

// ExtractUrlsFromRobots performs a URL extraction from robots by parsing it. This functions extracts ALL URLs, it doesn't respect ROBOTS.txt at all.
// In fact, it quite does the opposite
func ExtractUrlsFromRobots(ctx context.Context, bodyStr string, config *config.Config, qp *queue.Queue, domain string, registeredDomain string) {
	scanner := bufio.NewScanner(strings.NewReader(string(bodyStr)))

	for scanner.Scan() {
		line := scanner.Text()
		cleanLine := strings.TrimSpace(line)

		if cleanLine == "" || strings.HasPrefix(cleanLine, "#") {
			continue
		}

		cleanLine = strings.ToLower(cleanLine)
		if strings.HasPrefix(cleanLine, "disallow:") || strings.HasPrefix(cleanLine, "allow:") || strings.HasPrefix(cleanLine, "sitemap:") {
			values := strings.Split(cleanLine, ":")
			if len(values) == 1 {
				continue
			}

			values[1] = strings.TrimSpace(values[1])
			parsedUrl, err := urlvalidator.ValidateStringForUrl(config, values[1], domain, registeredDomain)
			if err == nil {
				qp.AddUrl(ctx, parsedUrl)
			}
		}
	}
}
