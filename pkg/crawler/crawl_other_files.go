package crawler

import (
	"context"
	"net/http"

	"github.com/gomills/gofocusedcrawler/internal/config"
	"github.com/gomills/gofocusedcrawler/pkg/queue"
	"github.com/gomills/gofocusedcrawler/pkg/urlextraction"
)

func CrawlOthers(ctx context.Context, config *config.Config, response *http.Response, qp *queue.Queue, domain string, registeredDomain string) {
	urlextraction.RegexFileForUrls(ctx, config, response, qp, domain, registeredDomain)
}
