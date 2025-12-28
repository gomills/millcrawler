package urlsextractors

import (
	"context"
	"net/url"
	"strings"

	http "github.com/bogdanfinn/fhttp"
	"github.com/gomills/gofocusedcrawler/internal/config"
	"github.com/gomills/gofocusedcrawler/pkg/queue"
)

// CrawlXml regexes sitemap.xml for urls.
func CrawlXml(ctx context.Context, config *config.Config, parsedUrl *url.URL, response *http.Response, qp *queue.Queue, domain string, registeredDomain string) {

	trimmedPath := strings.Trim(parsedUrl.Path, "/")

	if strings.HasSuffix(trimmedPath, "sitemap.xml") {
		regexFileForUrls(ctx, config, response, qp, domain, registeredDomain)

	}

}
