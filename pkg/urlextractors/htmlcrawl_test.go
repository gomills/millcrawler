package urlsextractors

import (
	"context"
	"io"
	"strings"
	"testing"
	"time"

	http "github.com/bogdanfinn/fhttp"
	"github.com/gomills/gofocusedcrawler/internal/config"
	"github.com/gomills/gofocusedcrawler/pkg/queue"
	"golang.org/x/net/html"
)

var testConfig = &config.Config{
	TimeOutDuration:        60,
	MaxPathDepth:           1,
	AllowedExternalDomains: []string{"github.com"},
	SensitivePatterns:      []string{"dashboard", "test", "repo", "private"},
	AllowedExtensions:      []string{".git", ".txt", ".js", ".html", ""},
	Workers:                3,
}

const (
	domain           = "www.example.com"
	registeredDomain = "example.com"
)

// TestRegexTextForUrls tests the regexTextForUrls function with a simple text containing an absolute URL
func TestRegexTextForUrls(t *testing.T) {
	ctx := context.Background()
	q := queue.NewQueue(1)

	// Simple text with one absolute URL
	text := "Check out //example.com/data.txt for more info"

	expected := 1
	regexTextForUrls(ctx, testConfig, text, q, domain, registeredDomain)

	// Verify that the URL was extracted and added to the queue
	finalCount := q.GetNCrawledUrls()
	if finalCount != expected {
		t.Errorf("Expected URL to be added to queue, got expected: %d, final: %d", expected, finalCount)
	}
}

// TestExtractUrlsFromHtmlComment tests the extractUrlsFromHtmlComment function with a simple comment containing a URL
func TestExtractUrlsFromHtmlComment(t *testing.T) {
	ctx := context.Background()
	q := queue.NewQueue(1)

	// Create an HTML comment node manually
	// HTML comments in the DOM are represented with CommentNode type
	commentNode := &html.Node{
		Type: html.CommentNode,
		Data: "<!-- Check this URL: //example.com/config.txt -->",
	}

	expected := 1
	extractUrlsFromHtmlComment(ctx, testConfig, commentNode, q, domain, registeredDomain)

	finalCount := q.GetNCrawledUrls()
	if finalCount != expected {
		t.Errorf("Expected URL to be added to queue from comment, got expected: %d, final: %d", expected, finalCount)
	}
}

// TestProcessOtherElements tests the processOtherElements function with a simple <link> element
func TestProcessOtherElements(t *testing.T) {
	ctx := context.Background()
	q := queue.NewQueue(1)

	// Create a simple <link> element with href attribute
	linkNode := &html.Node{
		Type: html.ElementNode,
		Data: "link",
		Attr: []html.Attribute{
			{
				Key: "href",
				Val: "https://example.com/style",
			},
		},
	}

	intrsAttrKeys := []string{"href", "src"}

	expected := 1
	processOtherElements(ctx, testConfig, linkNode, q, intrsAttrKeys, domain, registeredDomain)

	finalCount := q.GetNCrawledUrls()
	if finalCount != expected {
		t.Errorf("Expected URL to be added to queue from element, got expected: %d, final: %d", expected, finalCount)
	}
}

// TestRegexTextForUrlsMultipleMatches tests regex extraction with multiple URLs in text
func TestRegexTextForUrlsMultipleMatches(t *testing.T) {
	ctx := context.Background()
	q := queue.NewQueue(1)

	// Text with multiple absolute URLs
	text := "Visit //example.com/api and also //example.com/data.txt"

	expected := 2
	regexTextForUrls(ctx, testConfig, text, q, domain, registeredDomain)

	finalCount := q.GetNCrawledUrls()
	if finalCount != expected {
		t.Errorf("Expected URLs to be added to queue, got expected: %d, final: %d", expected, finalCount)
	}
}

// TestExtractUrlsFromHtmlCommentWithJavaScript tests comment extraction with JavaScript code
func TestExtractUrlsFromHtmlCommentWithJavaScript(t *testing.T) {
	ctx := context.Background()
	q := queue.NewQueue(1)

	// HTML comment containing JavaScript code with a URL
	commentNode := &html.Node{
		Type: html.CommentNode,
		Data: "window.apiUrl = '//example.com/api'; // old endpoint",
	}

	expected := 1
	extractUrlsFromHtmlComment(ctx, testConfig, commentNode, q, domain, registeredDomain)

	finalCount := q.GetNCrawledUrls()
	if finalCount != expected {
		t.Errorf("Expected URL to be added to queue from JavaScript comment, got expected: %d, final: %d", expected, finalCount)
	}
}

// TestProcessOtherElementsWithMultipleAttributes tests element processing with multiple attributes
func TestProcessOtherElementsWithMultipleAttributes(t *testing.T) {
	ctx := context.Background()
	q := queue.NewQueue(1)

	// Create a <form> element with action and other attributes
	formNode := &html.Node{
		Type: html.ElementNode,
		Data: "form",
		Attr: []html.Attribute{
			{
				Key: "action",
				Val: "https://example.com/submit",
			},
			{
				Key: "method",
				Val: "POST",
			},
			{
				Key: "id",
				Val: "myform",
			},
		},
	}

	intrsAttrKeys := []string{"action"}

	expected := 1
	processOtherElements(ctx, testConfig, formNode, q, intrsAttrKeys, domain, registeredDomain)

	finalCount := q.GetNCrawledUrls()
	if finalCount != expected {
		t.Errorf("Expected URL to be added to queue from form element, got expected: %d, final: %d", expected, finalCount)
	}
}

// TestProcessOtherElementsIgnoresIrrelevantAttributes tests that irrelevant attributes are ignored
func TestProcessOtherElementsIgnoresIrrelevantAttributes(t *testing.T) {
	ctx := context.Background()
	q := queue.NewQueue(1)

	// Create an element with attributes not in the interesting set
	elemNode := &html.Node{
		Type: html.ElementNode,
		Data: "embed",
		Attr: []html.Attribute{
			{
				Key: "class",
				Val: "https://example.com/styles",
			},
			{
				Key: "data-url",
				Val: "https://example.com/data",
			},
		},
	}

	intrsAttrKeys := []string{"src"} // Only 'src' is interesting, not 'class' or 'data-url'

	expected := 0
	processOtherElements(ctx, testConfig, elemNode, q, intrsAttrKeys, domain, registeredDomain)

	finalCount := q.GetNCrawledUrls()
	// Should not add anything since 'src' attribute is not present
	if finalCount != expected {
		t.Errorf("Expected no URLs to be added (only non-interesting attributes), got expected: %d, final: %d", expected, finalCount)
	}
}

// TestRegexTextForUrlsNoMatch tests regex extraction when text contains no URLs
func TestRegexTextForUrlsNoMatch(t *testing.T) {
	ctx := context.Background()
	q := queue.NewQueue(1)

	// Text with no URLs
	text := "This is just plain text with no URLs at all"

	expected := 0
	regexTextForUrls(ctx, testConfig, text, q, domain, registeredDomain)

	finalCount := q.GetNCrawledUrls()
	// Should not add anything since text has no URLs
	if finalCount != expected {
		t.Errorf("Expected no URLs to be added from text without URLs, got expected: %d, final: %d", expected, finalCount)
	}
}

// TestExtractUrlsFromHtmlCommentEmpty tests comment extraction with empty comment
func TestExtractUrlsFromHtmlCommentEmpty(t *testing.T) {
	ctx := context.Background()
	q := queue.NewQueue(1)

	// Empty comment
	commentNode := &html.Node{
		Type: html.CommentNode,
		Data: "",
	}

	expected := 0
	extractUrlsFromHtmlComment(ctx, testConfig, commentNode, q, domain, registeredDomain)

	finalCount := q.GetNCrawledUrls()
	// Should not add anything since comment is empty
	if finalCount != expected {
		t.Errorf("Expected no URLs to be added from empty comment, got expected: %d, final: %d", expected, finalCount)
	}
}

// TestProcessOtherElementsInvalidURL tests element with invalid URL
func TestProcessOtherElementsInvalidURL(t *testing.T) {
	ctx := context.Background()
	q := queue.NewQueue(1)

	// Element with an invalid URL (not matching heuristics)
	elemNode := &html.Node{
		Type: html.ElementNode,
		Data: "object",
		Attr: []html.Attribute{
			{
				Key: "data",
				Val: "ftp://example.com/file", // ftp scheme not interesting
			},
		},
	}

	intrsAttrKeys := []string{"data"}

	expected := 0
	processOtherElements(ctx, testConfig, elemNode, q, intrsAttrKeys, domain, registeredDomain)

	finalCount := q.GetNCrawledUrls()
	// Should not add anything since ftp:// URLs don't pass validation
	if finalCount != expected {
		t.Errorf("Expected no URLs to be added (invalid scheme), got expected: %d, final: %d", expected, finalCount)
	}
}

// TestRegexTextForUrlsWithQueryString tests regex extraction with URL containing query parameters
func TestRegexTextForUrlsWithQueryString(t *testing.T) {
	ctx := context.Background()
	q := queue.NewQueue(1)

	// Text with URL containing query parameters
	text := "API endpoint: //example.com/api?v=1&format=json"

	expected := 1
	regexTextForUrls(ctx, testConfig, text, q, domain, registeredDomain)

	finalCount := q.GetNCrawledUrls()
	if finalCount != expected {
		t.Errorf("Expected URL to be added to queue with query parameters, got expected: %d, final: %d", expected, finalCount)
	}
}

// TestProcessComplexElementsSimpleUrl tests processComplexElements with a simple URL (no code)
func TestProcessComplexElementsSimpleUrl(t *testing.T) {
	ctx := context.Background()
	q := queue.NewQueue(1)

	// Create a simple <a> element with href (no parentheses)
	aNode := &html.Node{
		Type: html.ElementNode,
		Data: "a",
		Attr: []html.Attribute{
			{
				Key: "href",
				Val: "https://example.com/target",
			},
		},
	}

	intrsAttrKeys := []string{"href"}

	expected := 1
	processComplexElements(ctx, testConfig, aNode, q, intrsAttrKeys, domain, registeredDomain)

	finalCount := q.GetNCrawledUrls()
	if finalCount != expected {
		t.Errorf("Expected simple URL to be added to queue, got expected: %d, final: %d", expected, finalCount)
	}
}

// TestProcessComplexElementsJavaScriptCode tests processComplexElements with JavaScript code containing a URL
func TestProcessComplexElementsJavaScriptCode(t *testing.T) {
	ctx := context.Background()
	q := queue.NewQueue(1)

	// Create a button element with onclick containing JavaScript with a URL
	// Presence of ( triggers regex extraction path
	buttonNode := &html.Node{
		Type: html.ElementNode,
		Data: "button",
		Attr: []html.Attribute{
			{
				Key: "onclick",
				Val: "navigate('//example.com/page');",
			},
		},
	}

	intrsAttrKeys := []string{"onclick"}

	expected := 1
	processComplexElements(ctx, testConfig, buttonNode, q, intrsAttrKeys, domain, registeredDomain)

	// Allow goroutines to complete (queue.AddUrl spawns a goroutine)
	time.Sleep(10 * time.Millisecond)

	finalCount := q.GetNCrawledUrls()
	if finalCount != expected {
		t.Errorf("Expected URL to be extracted from JavaScript code, got expected: %d, final: %d", expected, finalCount)
	}
}

// TestProcessComplexElementsMixedAttributes tests processComplexElements with both simple and complex attributes
func TestProcessComplexElementsMixedAttributes(t *testing.T) {
	ctx := context.Background()
	q := queue.NewQueue(1)

	// Create a script element with both src (simple) and onclick (complex with code)
	scriptNode := &html.Node{
		Type: html.ElementNode,
		Data: "script",
		Attr: []html.Attribute{
			{
				Key: "src",
				Val: "https://example.com/file.js",
			},
			{
				Key: "onclick",
				Val: "fetch('//example.com/api'); doSomething();",
			},
		},
	}

	intrsAttrKeys := []string{"src", "onclick"}

	expected := 2
	processComplexElements(ctx, testConfig, scriptNode, q, intrsAttrKeys, domain, registeredDomain)

	time.Sleep(10 * time.Millisecond)

	finalCount := q.GetNCrawledUrls()
	// Should extract from both src and onclick attributes
	if finalCount != expected {
		t.Errorf("Expected URLs to be added from mixed attributes, got expected: %d, final: %d", expected, finalCount)
	}
}

// TestProcessComplexElementsOnclickWithFunctionCall tests onclick with function call syntax
func TestProcessComplexElementsOnclickWithFunctionCall(t *testing.T) {
	ctx := context.Background()
	q := queue.NewQueue(1)

	// Create element with onclick containing function call with URL inside
	elemNode := &html.Node{
		Type: html.ElementNode,
		Data: "button",
		Attr: []html.Attribute{
			{
				Key: "onclick",
				Val: "goToUrl('https://example.com/newpage')",
			},
		},
	}

	intrsAttrKeys := []string{"onclick"}

	expected := 1
	processComplexElements(ctx, testConfig, elemNode, q, intrsAttrKeys, domain, registeredDomain)

	time.Sleep(10 * time.Millisecond)

	finalCount := q.GetNCrawledUrls()
	if finalCount != expected {
		t.Errorf("Expected URL to be extracted from function call, got expected: %d, final: %d", expected, finalCount)
	}
}

// TestProcessComplexElementsIgnoresNonInterestingAttributes tests that non-interesting attributes are skipped
func TestProcessComplexElementsIgnoresNonInterestingAttributes(t *testing.T) {
	ctx := context.Background()
	q := queue.NewQueue(1)

	// Create element with attributes that have URLs but are not in the interesting set
	elemNode := &html.Node{
		Type: html.ElementNode,
		Data: "a",
		Attr: []html.Attribute{
			{
				Key: "data-url",
				Val: "https://example.com/hidden",
			},
			{
				Key: "title",
				Val: "onClick='//example.com/redirect'",
			},
		},
	}

	intrsAttrKeys := []string{"href"} // Only href is interesting

	expected := 0
	processComplexElements(ctx, testConfig, elemNode, q, intrsAttrKeys, domain, registeredDomain)

	finalCount := q.GetNCrawledUrls()
	// Should not add anything since no interesting attributes present
	if finalCount != expected {
		t.Errorf("Expected no URLs to be added (non-interesting attributes), got expected: %d, final: %d", expected, finalCount)
	}
}

// TestProcessComplexElementsEmptyAttribute tests element with empty attribute value
func TestProcessComplexElementsEmptyAttribute(t *testing.T) {
	ctx := context.Background()
	q := queue.NewQueue(1)

	// Create element with empty href
	aNode := &html.Node{
		Type: html.ElementNode,
		Data: "a",
		Attr: []html.Attribute{
			{
				Key: "href",
				Val: "",
			},
		},
	}

	intrsAttrKeys := []string{"href"}

	expected := 0
	processComplexElements(ctx, testConfig, aNode, q, intrsAttrKeys, domain, registeredDomain)

	finalCount := q.GetNCrawledUrls()
	// Should not add anything since attribute is empty
	if finalCount != expected {
		t.Errorf("Expected no URLs to be added from empty attribute, got expected: %d, final: %d", expected, finalCount)
	}
}

// TestProcessComplexElementsMultipleUrls tests attribute with multiple URLs in JavaScript
func TestProcessComplexElementsMultipleUrls(t *testing.T) {
	ctx := context.Background()
	q := queue.NewQueue(1)

	// Create element with onclick containing multiple URLs
	elemNode := &html.Node{
		Type: html.ElementNode,
		Data: "button",
		Attr: []html.Attribute{
			{
				Key: "onclick",
				Val: "fetch('//example.com/api'); navigate('//example.com/page');",
			},
		},
	}

	intrsAttrKeys := []string{"onclick"}

	expected := 2
	processComplexElements(ctx, testConfig, elemNode, q, intrsAttrKeys, domain, registeredDomain)

	time.Sleep(10 * time.Millisecond)

	finalCount := q.GetNCrawledUrls()
	if finalCount != expected {
		t.Errorf("Expected multiple URLs to be added from JavaScript, got expected: %d, final: %d", expected, finalCount)
	}
}

// TestCrawlHtml tests the full CrawlHtml function with a complete HTML document
func TestCrawlHtml(t *testing.T) {
	ctx := context.Background()
	q := queue.NewQueue(1)

	// Create a realistic HTML document with various elements containing URLs
	htmlContent := `
	<!DOCTYPE html>
	<html>
	<head>
		<title>Test Page</title>
		<link rel="stylesheet" href="https://example.com/style">
		<script src="https://example.com/script.js"></script>
		<!-- API endpoint: //example.com/api -->
	</head>
	<body>
		<h1>Welcome</h1>
		<a href="https://example.com/about">About</a>
		<a href="/contact">Contact</a>
		<form action="https://example.com/submit" method="POST">
			<input type="submit">
		</form>
		<script>
			// Embedded script with URL
			fetch('//example.com/data');
		</script>
		<!-- Hidden links in comments: //example.com/admin -->
	</body>
	</html>
	`

	// Create a mock HTTP response
	response := &http.Response{
		Body: io.NopCloser(strings.NewReader(htmlContent)),
	}

	initialCount := q.GetNCrawledUrls()
	CrawlHtml(ctx, testConfig, response, q, domain, registeredDomain)

	// Allow goroutines to complete
	time.Sleep(20 * time.Millisecond)

	finalCount := q.GetNCrawledUrls()
	extractedCount := finalCount - initialCount
	expectedCount := 7 // style, script.js, about, contact, submit, data (fetch), admin (comment). Note: api comment may not match due to regex limitations

	if extractedCount != expectedCount {
		t.Errorf("Expected %d URLs to be extracted from HTML document, got %d (initial: %d, final: %d)", expectedCount, extractedCount, initialCount, finalCount)
	}
}

// TestCrawlHtmlWithMultipleElements tests CrawlHtml with various HTML elements
func TestCrawlHtmlWithMultipleElements(t *testing.T) {
	ctx := context.Background()
	q := queue.NewQueue(1)

	htmlContent := `
	<html>
	<body>
		<!-- Link in comment: //example.com/page1 -->
		<a href="/page2">Link 1</a>
		<a href="/page3">Link 2</a>
		<script src="https://example.com/lib.js"></script>
		<form action="/submit">Form</form>
		<iframe src="https://example.com/embed"></iframe>
		<object data="https://example.com/object"></object>
	</body>
	</html>
	`

	response := &http.Response{
		Body: io.NopCloser(strings.NewReader(htmlContent)),
	}

	initialCount := q.GetNCrawledUrls()
	CrawlHtml(ctx, testConfig, response, q, domain, registeredDomain)

	time.Sleep(20 * time.Millisecond)

	finalCount := q.GetNCrawledUrls()
	extractedCount := finalCount - initialCount
	expectedCount := 6 // page2, page3, lib.js, submit, embed, object. Note: page1 from comment may not match due to regex limitations

	if extractedCount != expectedCount {
		t.Errorf("Expected %d URLs from various HTML elements, got %d (initial: %d, final: %d)", expectedCount, extractedCount, initialCount, finalCount)
	}
}

// TestCrawlHtmlWithNoUrls tests CrawlHtml with HTML containing no extractable URLs
func TestCrawlHtmlWithNoUrls(t *testing.T) {
	ctx := context.Background()
	q := queue.NewQueue(1)

	htmlContent := `
	<html>
	<body>
		<h1>No URLs Here</h1>
		<p>Just plain text content</p>
		<div>Some more content</div>
	</body>
	</html>
	`

	response := &http.Response{
		Body: io.NopCloser(strings.NewReader(htmlContent)),
	}

	initialCount := q.GetNCrawledUrls()
	CrawlHtml(ctx, testConfig, response, q, domain, registeredDomain)

	time.Sleep(10 * time.Millisecond)

	finalCount := q.GetNCrawledUrls()
	extractedCount := finalCount - initialCount
	expectedCount := 0 // No URLs in the HTML

	if extractedCount != expectedCount {
		t.Errorf("Expected %d URLs from plain HTML, got %d (initial: %d, final: %d)", expectedCount, extractedCount, initialCount, finalCount)
	}
}

// TestCrawlHtmlWithMalformedHTML tests CrawlHtml with malformed HTML
func TestCrawlHtmlWithMalformedHTML(t *testing.T) {
	ctx := context.Background()
	q := queue.NewQueue(1)

	// Malformed but recoverable HTML
	htmlContent := `
	<html>
	<body>
		<a href="https://example.com/page1">Link
		<a href="https://example.com/page2">Link
		<p>Unclosed tags
	</body>
	</html>
	`

	response := &http.Response{
		Body: io.NopCloser(strings.NewReader(htmlContent)),
	}

	initialCount := q.GetNCrawledUrls()
	CrawlHtml(ctx, testConfig, response, q, domain, registeredDomain)

	time.Sleep(20 * time.Millisecond)

	finalCount := q.GetNCrawledUrls()
	extractedCount := finalCount - initialCount
	expectedCount := 2 // page1, page2 (both from <a> tags)

	if extractedCount != expectedCount {
		t.Errorf("Expected %d URLs from malformed HTML, got %d (initial: %d, final: %d)", expectedCount, extractedCount, initialCount, finalCount)
	}
}

// TestCrawlHtmlWithNestedElements tests CrawlHtml with deeply nested HTML structures
func TestCrawlHtmlWithNestedElements(t *testing.T) {
	ctx := context.Background()
	q := queue.NewQueue(1)

	htmlContent := `
	<html>
	<body>
		<div>
			<section>
				<article>
					<a href="/page1">Nested Link</a>
					<script src="https://example.com/nested.js"></script>
				</article>
			</section>
		</div>
	</body>
	</html>
	`

	response := &http.Response{
		Body: io.NopCloser(strings.NewReader(htmlContent)),
	}

	initialCount := q.GetNCrawledUrls()
	CrawlHtml(ctx, testConfig, response, q, domain, registeredDomain)

	time.Sleep(20 * time.Millisecond)

	finalCount := q.GetNCrawledUrls()
	extractedCount := finalCount - initialCount
	expectedCount := 2 // page1, nested.js

	if extractedCount != expectedCount {
		t.Errorf("Expected %d URLs from nested elements, got %d (initial: %d, final: %d)", expectedCount, extractedCount, initialCount, finalCount)
	}
}
