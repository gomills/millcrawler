package crawler

import (
	"context"
	"net/http"

	"github.com/gomills/gofocusedcrawler/internal/config"
	"github.com/gomills/gofocusedcrawler/pkg/queue"
	"github.com/gomills/gofocusedcrawler/pkg/urlextraction"
	"golang.org/x/net/html"
)

func CrawlHtml(ctx context.Context, config *config.Config, response *http.Response, qp *queue.Queue, domain string, registeredDomain string) {
	rootNode, err := html.Parse(response.Body)
	if err != nil {
		return
	}

	urlextraction.ExtractUrlsFromHtmlNode(ctx, config, rootNode, qp, domain, registeredDomain)
}
