package urls_extractors

import (
	"context"

	http "github.com/bogdanfinn/fhttp"
	"github.com/gomills/gofocusedcrawler/internal/config"
	"github.com/gomills/gofocusedcrawler/pkg/queue"
)

// CrawlOthers regexes the file for urls. Only called on error payloads.
func CrawlOthers(ctx context.Context, config *config.Config, response *http.Response, qp *queue.Queue, domain string, registeredDomain string) {
	regexFileForUrls(ctx, config, response, qp, domain, registeredDomain)
}
