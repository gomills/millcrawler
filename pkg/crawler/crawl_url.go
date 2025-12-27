package crawler

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"path/filepath"

	http "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
	"github.com/gomills/gofocusedcrawler/internal/config"
	"github.com/gomills/gofocusedcrawler/pkg/queue"
	urls_extr "github.com/gomills/gofocusedcrawler/pkg/urls_extractors"
	"github.com/gomills/gofocusedcrawler/pkg/utils"
)

// crawlUrl returns error only on 429 status_code hit
func crawlUrl(ctx context.Context, client tls_client.HttpClient, headers http.Header, config *config.Config, url *url.URL, qp *queue.Queue, domain string, registeredDomain string) error {

	// skip crawling urls from external domains
	isLocal, err := utils.IsLocalUrl(url, registeredDomain)
	if err != nil || !isLocal {
		return nil
	}

	// craft request with context
	request, _ := http.NewRequestWithContext(ctx, "GET", url.String(), nil)
	request.Header = headers

	// send it
	response, err := client.Do(request)
	if err != nil {
		// log.Printf("%s failed request for %s", id, url.String())
		return nil
	}
	defer response.Body.Close()

	// TODO: crawl headers for all cases

	// check status code.
	// For 429 cancel context.
	// For 2xx and others proceed to crawling.
	// log.Printf("%s %d request for %s", id, response.StatusCode, url.String())
	if response.StatusCode > 299 {

		if response.StatusCode == 429 {
			return fmt.Errorf("429_hit")
		}

		// crawl error payload body with regex for debugging endpoints
		urls_extr.CrawlOthers(ctx, config, response, qp, domain, registeredDomain)

		return nil
	}

	ext := filepath.Ext(url.Path)

	// crawl the body according to the extension
	switch ext {

	case "", ".htm", ".html":
		urls_extr.CrawlHtml(ctx, config, response, qp, domain, registeredDomain)

	case ".js", ".min.js":
		CrawlJavascript(ctx, config, response, qp, domain, registeredDomain)

	case ".txt":
		urls_extr.CrawlTxt(ctx, config, url, response, qp, domain, registeredDomain)

	case ".xml":
		urls_extr.CrawlXml(ctx, config, url, response, qp, domain, registeredDomain)

	default:
		// urls_extr.CrawlOthers(ctx, config, response, qp, domain, registeredDomain)
	}

	return nil

}

// This function is here because the linter just won't find it when I move it to another file!
func CrawlJavascript(ctx context.Context, config *config.Config, response *http.Response, qp *queue.Queue, domain string, registeredDomain string) {
	jsCodeByte, err := io.ReadAll(response.Body)
	if err != nil {
		return
	}
	urls_extr.CrawlJs(ctx, config, jsCodeByte, qp, domain, registeredDomain)
}
