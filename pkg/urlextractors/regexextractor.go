package urlsextractors

import (
	"context"
	"io"
	"regexp"

	"github.com/gomills/millcrawler/pkg/config"
	"github.com/gomills/millcrawler/pkg/queue"
	"github.com/gomills/millcrawler/pkg/urlvalidator"
)

var AbsoluteUrlPattern = regexp.MustCompile(`//[a-zA-Z0-9.-]+(?:/[^\s.,:"')<]*(?:\.[^\s.,:"')<]+)*/?)`)

// regexFileForUrls takes a raw response and regexes it for urls.
func regexFileForUrls(ctx context.Context, config *config.Config, responseBody io.ReadCloser, qp *queue.Queue, domain string, registeredDomain string) {

	body, err := io.ReadAll(responseBody)
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
