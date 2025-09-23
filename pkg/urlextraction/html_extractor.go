package urlextraction

import (
	"context"
	"strings"

	"github.com/gomills/gofocusedcrawler/internal/config"
	"github.com/gomills/gofocusedcrawler/pkg/queue"
	"github.com/gomills/gofocusedcrawler/pkg/urlvalidator"
	"golang.org/x/net/html"
)

// elementInterestingAttributes maps HTML elements to their attributes that may contain URLs. This is completely heuristic.
var elementInterestingAttributes = map[string][]string{
	"a":      {"href"},
	"script": {"src", "content", "href", "onclick", "action", "formaction", "codebase"},
	"button": {"onclick", "formaction"},
	"link":   {"href", "src"},
	"form":   {"action"},
	"object": {"data"},
	"embed":  {"src"},
	"iframe": {"src"},
}

// ExtractUrlsFromHtmlNode recursively traverses an HTML node tree, extracting URLs from elements and comments.
func ExtractUrlsFromHtmlNode(ctx context.Context, config *config.Config, node *html.Node, q *queue.Queue, domain string, registeredDomain string) {

	if node.Type == html.ElementNode {
		elementKey := node.Data
		if intersAttrKeys, ok := elementInterestingAttributes[elementKey]; ok {
			extractUrlsFromHtmlElement(ctx, config, node, q, intersAttrKeys, domain, registeredDomain)
		}
	}

	if node.Type == html.CommentNode {
		extractUrlsFromHtmlComment(ctx, config, node, q, domain, registeredDomain)
	}

	for i := node.FirstChild; i != nil; i = i.NextSibling {
		ExtractUrlsFromHtmlNode(ctx, config, i, q, domain, registeredDomain)
	}

}

// ExtractUrlsFromElement extracts URLs from the specified attributes of an HTML element node.
func extractUrlsFromHtmlElement(ctx context.Context, config *config.Config, n *html.Node, q *queue.Queue, intrsAttrKeys []string, domain string, registeredDomain string) {

	if n.Data == "a" || n.Data == "script" || n.Data == "button" {
		processComplexElements(ctx, config, n, q, intrsAttrKeys, domain, registeredDomain)

	} else {
		processOtherElements(ctx, config, n, q, intrsAttrKeys, domain, registeredDomain)
	}

}

// processComplexElements handles extraction for elements with potentially complex attribute values (e.g., JavaScript in onclick).
func processComplexElements(ctx context.Context, config *config.Config, n *html.Node, q *queue.Queue, intrsAttrKeys []string, domain string, registeredDomain string) {
	// Iterate over all attributes of the node

	for _, att := range n.Attr {

		for _, attrKey := range intrsAttrKeys {

			if attrKey == att.Key {

				// If the attribute value looks like it contains code, extract URLs using regex
				if x := strings.IndexByte(att.Val, '('); x != -1 {
					regexTextForUrls(ctx, config, att.Val, q, domain, registeredDomain)
				} else {
					// Otherwise, validate the URL directly
					parsedUrl, err := urlvalidator.ValidateStringForUrl(config, att.Val, domain, registeredDomain)
					if err == nil {
						q.AddUrl(ctx, parsedUrl)
					}
				}

			}
		}
	}
}

// processOtherElements handles extraction for elements with simple attribute values.
func processOtherElements(ctx context.Context, config *config.Config, n *html.Node, q *queue.Queue, intrsAttrKeys []string, domain string, registeredDomain string) {
	// Iterate over all attributes of the node

	for _, att := range n.Attr {

		for _, attrKey := range intrsAttrKeys {

			if attrKey == att.Key {

				// Validate and add the URL if valid
				parsedUrl, err := urlvalidator.ValidateStringForUrl(config, att.Val, domain, registeredDomain)
				if err == nil {
					q.AddUrl(ctx, parsedUrl)
				}

			}
		}
	}
}

// ExtractUrlsFromHtmlComment extracts URLs from HTML comment nodes using regex.
func extractUrlsFromHtmlComment(ctx context.Context, config *config.Config, n *html.Node, q *queue.Queue, domain string, registeredDomain string) {
	regexTextForUrls(ctx, config, n.Data, q, domain, registeredDomain)
}

// regexTextForUrls finds all absolute URLs in the given text using a regex pattern and adds them to the queue if valid.
func regexTextForUrls(ctx context.Context, config *config.Config, text string, q *queue.Queue, domain string, registeredDomain string) {

	if matches := AbsoluteUrlPattern.FindAllString(text, -1); matches != nil {

		for _, match := range matches {

			parsedUrl, err := urlvalidator.ValidateStringForUrl(config, match, domain, registeredDomain)
			if err == nil {
				q.AddUrl(ctx, parsedUrl)
			}

		}
	}
}
