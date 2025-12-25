package crawler

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"path/filepath"

	"github.com/gomills/gofocusedcrawler/internal/config"
	"github.com/gomills/gofocusedcrawler/pkg/queue"
	urls_extr "github.com/gomills/gofocusedcrawler/pkg/urls_extractors"
	"github.com/gomills/gofocusedcrawler/pkg/utils"
)

// crawlUrl returns true
func crawlUrl(ctx context.Context, config *config.Config, url *url.URL, qp *queue.Queue, domain string, registeredDomain string, id string) bool {

	// Skip crawling external domain URLs
	isLocal, err := utils.IsLocalUrl(url, registeredDomain)
	if err != nil || !isLocal {
		return false
	}

	// Craft request with context
	request, _ := http.NewRequestWithContext(ctx, "GET", url.String(), nil)
	request.Header.Set("User-Agent", config.CustomUserAgent)

	// Send request
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		// log.Printf("%s failed request for %s", id, url.String())
		return false
	}
	defer response.Body.Close()

	// Check status code. If 429 cancel context to stop whole crawler. If other than 2XX, try to scrap some URL anyway with regex
	// log.Printf("%s %d request for %s", id, response.StatusCode, url.String())
	if response.StatusCode > 299 {

		if response.StatusCode == 429 {
			return true
		}

		urls_extr.CrawlOthers(ctx, config, response, qp, domain, registeredDomain)

		return false
	}

	ext := filepath.Ext(url.Path)

	// Crawl the body according to the extension
	switch ext {
	case "", ".htm", ".html":
		urls_extr.CrawlHtml(ctx, config, response, qp, domain, registeredDomain)
	case ".js", ".min.js":
		CrawlJavascript(ctx, config, response, qp, domain, registeredDomain)
	case ".txt":
		urls_extr.CrawlTxtFileUrl(ctx, config, url, response, qp, domain, registeredDomain)
	default:
		urls_extr.CrawlOthers(ctx, config, response, qp, domain, registeredDomain)
	}

	return false

}

// This function is here because the linter just won't find it when I move it to another file!
func CrawlJavascript(ctx context.Context, config *config.Config, response *http.Response, qp *queue.Queue, domain string, registeredDomain string) {
	jsCodeByte, err := io.ReadAll(response.Body)
	if err != nil {
		return
	}
	urls_extr.TraverseJs(ctx, config, jsCodeByte, qp, domain, registeredDomain)
}
