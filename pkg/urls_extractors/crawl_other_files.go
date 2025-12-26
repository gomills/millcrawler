package urls_extractors

import (
	"context"
	"net/http"

	"github.com/gomills/gofocusedcrawler/internal/config"
	"github.com/gomills/gofocusedcrawler/pkg/queue"
)

// CrawlOthers regexes the file for urls. Danger for this, we'll limit its usage soon.
func CrawlOthers(ctx context.Context, config *config.Config, response *http.Response, qp *queue.Queue, domain string, registeredDomain string) {
	regexFileForUrls(ctx, config, response, qp, domain, registeredDomain)
}
