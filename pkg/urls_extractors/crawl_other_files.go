package urls_extractors

import (
	"context"
	"net/http"

	"github.com/gomills/gofocusedcrawler/internal/config"
	"github.com/gomills/gofocusedcrawler/pkg/queue"
)

func CrawlOthers(ctx context.Context, config *config.Config, response *http.Response, qp *queue.Queue, domain string, registeredDomain string) {
	RegexFileForUrls(ctx, config, response, qp, domain, registeredDomain)
}
