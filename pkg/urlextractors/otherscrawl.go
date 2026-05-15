package urlsextractors

import (
	"context"
	"io"

	"github.com/gomills/millcrawler/pkg/config"
	"github.com/gomills/millcrawler/pkg/queue"
)

// CrawlOthers regexes the file for urls. Only called on error payloads.
func CrawlOthers(ctx context.Context, config *config.Config, responseBody io.ReadCloser, qp *queue.Queue, domain string, registeredDomain string) {
	regexFileForUrls(ctx, config, responseBody, qp, domain, registeredDomain)
}
