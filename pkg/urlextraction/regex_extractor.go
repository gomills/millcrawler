package urlextraction

import (
	"context"
	"io"
	"net/http"
	"regexp"

	"github.com/gomills/gofocusedcrawler/internal/config"
	"github.com/gomills/gofocusedcrawler/pkg/queue"
	"github.com/gomills/gofocusedcrawler/pkg/urlvalidator"
)

var AbsoluteUrlPattern = regexp.MustCompile(`//[a-zA-Z0-9.-]+(?:/[^\s.,:"')]*)?`)

func RegexFileForUrls(ctx context.Context, config *config.Config, response *http.Response, qp *queue.Queue, domain string, registeredDomain string) {

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return
	}

	if matches := AbsoluteUrlPattern.FindAll(body, -1); matches != nil {
		processUrlsMatches(ctx, config, matches, qp, domain, registeredDomain)
	}

}

func processUrlsMatches(ctx context.Context, config *config.Config, matches [][]byte, qp *queue.Queue, domain string, registeredDomain string) {

	for _, match := range matches {

		stringMatch := string(match)

		parsedUrl, err := urlvalidator.ValidateStringForUrl(config, stringMatch, domain, registeredDomain)
		if err == nil {
			qp.AddUrl(ctx, parsedUrl)
		}

	}
}
