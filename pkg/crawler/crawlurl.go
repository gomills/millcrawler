package crawler

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/url"
	"path/filepath"

	http "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
	"github.com/gomills/millcrawler/pkg/config"
	"github.com/gomills/millcrawler/pkg/queue"
	urls_extr "github.com/gomills/millcrawler/pkg/urlextractors"
	"github.com/gomills/millcrawler/pkg/utilities"
)

// crawlUrl requests a URL and tries to extract URLs and secrets(optional), adding both to the queue.
func crawlUrl(ctx context.Context, client tls_client.HttpClient, config *config.Config, url *url.URL, qp *queue.Queue, domain string, registeredDomain string) error {

	// 1.1. craft request with context
	request, err := http.NewRequestWithContext(ctx, "GET", url.String(), nil)
	if err != nil {
		return nil
	}

	// 1.2. send it
	response, err := client.Do(request)
	if err != nil {
		// log.Printf("%s failed request for %s", id, url.String())
		return nil
	}

	// 2.1. download body into memory
	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil
	}
	defer response.Body.Close()

	// 2.2. but this read leaves response.Body used. Parsers need a io.ReadCloser interface instead of []byte;
	// restore response's body so that parsers can use it
	response.Body = io.NopCloser(bytes.NewReader(bodyBytes))

	// 3.1. check status code.
	// For 2xx and others proceed to crawling.
	// For 429 return error
	// For other errors crawl payload with regex.
	// log.Printf("%s %d request for %s", id, response.StatusCode, url.String())
	if response.StatusCode > 299 {

		if response.StatusCode == 429 {
			return errors.New("antibot status code")
		}

		// 3.2. crawl error payloads with regex for debugging endpoints
		urls_extr.CrawlOthers(ctx, config, response.Body, qp, domain, registeredDomain)

		// 3.3. scan for secrets in headers + response's body in error payloads
		if config.ScanSecrets {
			scanSecrets(qp, response.Header, url.String(), bodyBytes)
		}

		return nil
	}

	// 4.0. from external domains don't crawl, just look for secrets if enabled
	isLocal, err := utilities.IsLocalUrl(url, registeredDomain)
	if err != nil || !isLocal {
		if config.ScanSecrets {
			scanSecrets(qp, response.Header, url.String(), bodyBytes)
		}
		return nil
	}

	ext := filepath.Ext(url.Path)

	// 4.1. scan for secrets except for .html which make up the bulk of urls and are not the most interesting
	if ext != "" && ext != ".htm" && ext != ".html" && config.ScanSecrets {
		scanSecrets(qp, response.Header, url.String(), bodyBytes)
	}

	// 4.2. crawl the body according to the extension
	switch ext {

	case "", ".htm", ".html":
		urls_extr.CrawlHtml(ctx, config, response.Body, qp, domain, registeredDomain)

	case ".js", ".min.js":
		urls_extr.CrawlJs(ctx, config, response.Body, qp, domain, registeredDomain)

	case ".txt":
		urls_extr.CrawlTxt(ctx, config, url, response.Body, qp, domain, registeredDomain)

	case ".xml":
		urls_extr.CrawlXml(ctx, config, url, response.Body, qp, domain, registeredDomain)

	default:
		// do nothing. A possible implementation would be to regex at a high cost of resources. For specific filetypes,
		// include a respective parser.
		// urls_extr.CrawlOthers(ctx, config, response, qp, domain, registeredDomain)
	}

	return nil

}
