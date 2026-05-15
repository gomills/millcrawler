package urlsextractors

import (
	"context"
	"io"
	"net/url"
	"strings"

	"github.com/gomills/millcrawler/pkg/config"
	"github.com/gomills/millcrawler/pkg/queue"
)

// CrawlTxtFileUrl parses robots.txt for urls.
func CrawlTxt(ctx context.Context, config *config.Config, parsedUrl *url.URL, robotsFile io.ReadCloser, qp *queue.Queue, domain string, registeredDomain string) {

	if strings.Contains(parsedUrl.Path, "robots.txt") {
		crawlRobots(ctx, config, robotsFile, qp, domain, registeredDomain)

	} else {
		// regexFileForUrls(ctx, config, response, qp, domain, registeredDomain)
	}

}

func crawlRobots(ctx context.Context, config *config.Config, robotsFile io.ReadCloser, qp *queue.Queue, domain string, registeredDomain string) {
	bodyStr, err := io.ReadAll(robotsFile)
	if err != nil {
		// log.Print("Failed to read robots.txt' body")
		return
	}

	extractUrlsFromRobots(ctx, string(bodyStr), config, qp, domain, registeredDomain)
}
