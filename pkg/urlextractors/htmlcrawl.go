package urlsextractors

import (
	"context"
	"strings"

	http "github.com/bogdanfinn/fhttp"
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
	"meta":   {"content"},                    // gpt
	"base":   {"href"},                       // gpt
	"input":  {"value", "formaction", "src"}, // gpt
}

// CrawlHtml iterates over an .html file to extract urls from element's attributes. Returns nothing because
// appends urls to the passed queue.
func CrawlHtml(ctx context.Context, config *config.Config, response *http.Response, qp *queue.Queue, domain string, registeredDomain string) {

	// get the root node
	rootNode, err := html.Parse(response.Body)
	if err != nil {
		return
	}

	// iterate over it to complete the tree
	extractUrlsFromHtmlNode(ctx, config, rootNode, qp, domain, registeredDomain)

}

// extractUrlsFromHtmlNode recursively traverses an HTML node tree, extracting URLs from elements and comments.
func extractUrlsFromHtmlNode(ctx context.Context, config *config.Config, node *html.Node, q *queue.Queue, domain string, registeredDomain string) {

	// extract urls from elements: a, href, etc.
	if node.Type == html.ElementNode {

		// get its key (the name, e.g href) and check it's in the interesting elements map
		elementKey := node.Data
		if intersAttrKeys, ok := elementInterestingAttributes[elementKey]; ok {
			extractUrlsFromHtmlElement(ctx, config, node, q, intersAttrKeys, domain, registeredDomain)
		}

	}

	// extract urls from .html comments <-- xxx --!>
	if node.Type == html.CommentNode {
		extractUrlsFromHtmlComment(ctx, config, node, q, domain, registeredDomain)
	}

	// recursively iterate children's siblings
	for i := node.FirstChild; i != nil; i = i.NextSibling {
		extractUrlsFromHtmlNode(ctx, config, i, q, domain, registeredDomain)
	}

}

// extractUrlsFromElement extracts URLs from the specified attributes of an HTML element node.
func extractUrlsFromHtmlElement(ctx context.Context, config *config.Config, n *html.Node, q *queue.Queue, intrsAttrKeys []string, domain string, registeredDomain string) {

	// we'll give these a special treatment for being the most common and vulnerable
	if n.Data == "a" || n.Data == "script" || n.Data == "button" {
		processComplexElements(ctx, config, n, q, intrsAttrKeys, domain, registeredDomain)

		// these won't have added complexity
	} else {
		processOtherElements(ctx, config, n, q, intrsAttrKeys, domain, registeredDomain)
	}

}

// processComplexElements handles extraction for elements with potentially complex attribute values (e.g., JavaScript in onclick).
func processComplexElements(ctx context.Context, config *config.Config, n *html.Node, q *queue.Queue, intrsAttrKeys []string, domain string, registeredDomain string) {

	// iterate over all attributes of the node
	for _, att := range n.Attr {

		// check it's in the allowed set
		for _, attrKey := range intrsAttrKeys {

			if attrKey == att.Key {

				// if the attribute value looks like it contains code, fallback to regex
				if x := strings.IndexByte(att.Val, '('); x != -1 {
					regexTextForUrls(ctx, config, att.Val, q, domain, registeredDomain)

				} else {

					// otherwise, validate the URL directly
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

	// iterate over all attributes of the node
	for _, att := range n.Attr {

		// check it's in the allowed set
		for _, attrKey := range intrsAttrKeys {

			if attrKey == att.Key {

				// validate and add the URL if valid
				parsedUrl, err := urlvalidator.ValidateStringForUrl(config, att.Val, domain, registeredDomain)
				if err == nil {
					q.AddUrl(ctx, parsedUrl)
				}

			}
		}
	}
}

// extractUrlsFromHtmlComment extracts URLs from HTML comment nodes using regex.
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
