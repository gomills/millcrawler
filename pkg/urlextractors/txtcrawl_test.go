package urlsextractors

import (
	"context"
	"io"
	"net/url"
	"slices"
	"strings"
	"testing"
	"time"

	http "github.com/bogdanfinn/fhttp"
	"github.com/gomills/gofocusedcrawler/internal/config"
	"github.com/gomills/gofocusedcrawler/pkg/queue"
)

var txtTestConfig = &config.Config{
	TimeOutDuration:        60,
	MaxPathDepth:           2,
	AllowedExternalDomains: []string{"cdn.example.com"},
	SensitivePatterns:      []string{"admin", "private"},
	AllowedExtensions:      []string{".txt", ".json", ".xml", ""},
	Workers:                3,
}

const (
	txtDomain           = "www.example.com"
	txtRegisteredDomain = "example.com"
)

// TestCrawlTxtRobotsFile tests that robots.txt URLs are correctly discriminated and extracted
func TestCrawlTxtRobotsFile(t *testing.T) {
	ctx := context.Background()
	q := queue.NewQueue(1)

	robotsContent := `
User-agent: *
Disallow: /api/v1/data
Allow: /public/search
Sitemap: //cdn.example.com/sitemap.xml
`

	robotsUrl, _ := url.Parse("https://www.example.com/robots.txt")
	response := &http.Response{
		Body: io.NopCloser(strings.NewReader(robotsContent)),
	}

	CrawlTxt(ctx, txtTestConfig, robotsUrl, response, q, txtDomain, txtRegisteredDomain)
	time.Sleep(10 * time.Millisecond)

	extractedUrls := q.GetFoundUrls()
	slices.Sort(extractedUrls)

	// Expected URLs extracted from robots.txt
	expectedUrls := []string{
		"https://www.example.com/public/search",
		"https://cdn.example.com/sitemap.xml",
	}
	slices.Sort(expectedUrls)

	if len(extractedUrls) != len(expectedUrls) {
		t.Logf("Extracted URLs (%d): %v", len(extractedUrls), extractedUrls)
		t.Logf("Expected URLs (%d): %v", len(expectedUrls), expectedUrls)
		t.Fatalf("Expected %d URLs, got %d", len(expectedUrls), len(extractedUrls))
	}

	for i, expected := range expectedUrls {
		if i >= len(extractedUrls) || extractedUrls[i] != expected {
			t.Errorf("URL mismatch at index %d: expected %q, got %q", i, expected, extractedUrls[i])
		}
	}
}

// TestCrawlTxtNonRobotsFile tests that non-robots.txt files are not processed for URL extraction
func TestCrawlTxtNonRobotsFile(t *testing.T) {
	ctx := context.Background()
	q := queue.NewQueue(1)

	textContent := `
Some random text with URLs:
Disallow: /api/v1/data
Allow: /public/search
Sitemap: //cdn.example.com/sitemap.xml
`

	nonRobotsUrl, _ := url.Parse("https://www.example.com/data.txt")
	response := &http.Response{
		Body: io.NopCloser(strings.NewReader(textContent)),
	}

	initialCount := q.GetNCrawledUrls()
	CrawlTxt(ctx, txtTestConfig, nonRobotsUrl, response, q, txtDomain, txtRegisteredDomain)
	time.Sleep(10 * time.Millisecond)

	finalCount := q.GetNCrawledUrls()
	extractedCount := finalCount - initialCount

	if extractedCount != 0 {
		extractedUrls := q.GetFoundUrls()
		t.Errorf("Expected 0 URLs from non-robots.txt file, got %d: %v", extractedCount, extractedUrls)
	}
}
