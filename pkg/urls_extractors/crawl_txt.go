package urls_extractors

import (
	"context"
	"io"
	"log"
	"net/url"
	"strings"

	http "github.com/bogdanfinn/fhttp"
	"github.com/gomills/gofocusedcrawler/internal/config"
	"github.com/gomills/gofocusedcrawler/pkg/queue"
)

// CrawlTxtFileUrl parses robots.txt for urls.
func CrawlTxt(ctx context.Context, config *config.Config, parsedUrl *url.URL, response *http.Response, qp *queue.Queue, domain string, registeredDomain string) {

	if strings.Contains(parsedUrl.Path, "robots.txt") {
		crawlRobots(ctx, config, response, qp, domain, registeredDomain)

	} else {
		// regexFileForUrls(ctx, config, response, qp, domain, registeredDomain)
	}

}

func crawlRobots(ctx context.Context, config *config.Config, response *http.Response, qp *queue.Queue, domain string, registeredDomain string) {
	bodyStr, err := io.ReadAll(response.Body)
	if err != nil {
		log.Print("Failed to read robots.txt' body")
		return
	}

	extractUrlsFromRobots(ctx, string(bodyStr), config, qp, domain, registeredDomain)
}
