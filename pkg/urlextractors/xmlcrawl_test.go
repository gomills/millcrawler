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
	"github.com/gomills/millcrawler/pkg/config"
	"github.com/gomills/millcrawler/pkg/queue"
)

var xmlTestConfig = &config.Config{
	TimeOutDuration:        60,
	MaxPathDepth:           2,
	AllowedExternalDomains: []string{},
	SensitivePatterns:      []string{},
	AllowedExtensions:      []string{".html", ".xml", ""},
	Workers:                3,
}

const (
	xmlDomain           = "www.example.com"
	xmlRegisteredDomain = "example.com"
)

// TestCrawlXmlSitemapExtraction tests that URLs are extracted from sitemap.xml files
func TestCrawlXmlSitemapExtraction(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	q := queue.NewQueue(1)

	// Mock sitemap.xml content with 2 URLs
	sitemapContent := `<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
  <url>
    <loc>//example.com/page1</loc>
    <lastmod>2024-01-01</lastmod>
  </url>
  <url>
    <loc>https://example.com/page2</loc>
    <lastmod>2024-01-02</lastmod>
  </url>
</urlset>`

	sitemapUrl, _ := url.Parse("https://www.example.com/sitemap.xml")
	response := &http.Response{
		Body: io.NopCloser(strings.NewReader(sitemapContent)),
	}

	CrawlXml(ctx, xmlTestConfig, sitemapUrl, response.Body, q, xmlDomain, xmlRegisteredDomain)
	time.Sleep(10 * time.Millisecond)

	extractedUrls := q.GetFoundUrls()
	slices.Sort(extractedUrls)

	// Expected URLs extracted from sitemap
	expectedUrls := []string{
		"https://example.com/page1",
		"https://example.com/page2",
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

// TestCrawlXmlNonSitemap tests that non-sitemap XML files are ignored
func TestCrawlXmlNonSitemap(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	q := queue.NewQueue(1)

	xmlContent := `<?xml version="1.0"?>
<root>
  <url>//example.com/page</url>
</root>`

	nonSitemapUrl, _ := url.Parse("https://www.example.com/data.xml")
	response := &http.Response{
		Body: io.NopCloser(strings.NewReader(xmlContent)),
	}

	CrawlXml(ctx, xmlTestConfig, nonSitemapUrl, response.Body, q, xmlDomain, xmlRegisteredDomain)
	time.Sleep(10 * time.Millisecond)

	// Non-sitemap.xml files should not be processed
	if count := q.GetNFoundUrls(); count != 0 {
		t.Errorf("Expected 0 URLs from non-sitemap XML file, got %d", count)
	}
}
