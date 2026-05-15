package urlsextractors

import (
	"context"
	"io"
	"net/url"
	"strings"

	"github.com/gomills/millcrawler/pkg/config"
	"github.com/gomills/millcrawler/pkg/queue"
)

// CrawlXml regexes sitemap.xml for urls.
func CrawlXml(ctx context.Context, config *config.Config, parsedUrl *url.URL, siteMapFile io.ReadCloser, qp *queue.Queue, domain string, registeredDomain string) {

	trimmedPath := strings.Trim(parsedUrl.Path, "/")

	if strings.HasSuffix(trimmedPath, "sitemap.xml") {
		regexFileForUrls(ctx, config, siteMapFile, qp, domain, registeredDomain)

	}

}
