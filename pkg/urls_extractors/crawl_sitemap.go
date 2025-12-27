package urls_extractors

import (
	"context"
	"net/http"
	"net/url"
	"strings"

	"github.com/gomills/gofocusedcrawler/internal/config"
	"github.com/gomills/gofocusedcrawler/pkg/queue"
)

// CrawlXml regexes sitemap.xml for urls.
func CrawlXml(ctx context.Context, config *config.Config, parsedUrl *url.URL, response *http.Response, qp *queue.Queue, domain string, registeredDomain string) {

	if strings.Contains(parsedUrl.Path, "sitemap.xml") {
		regexFileForUrls(ctx, config, response, qp, domain, registeredDomain)

	}

}
