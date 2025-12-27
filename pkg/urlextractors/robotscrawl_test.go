package urlsextractors

import (
	"context"
	"slices"
	"testing"
	"time"

	"github.com/gomills/gofocusedcrawler/internal/config"
	"github.com/gomills/gofocusedcrawler/pkg/queue"
)

var robotsTestConfig = &config.Config{
	TimeOutDuration:        60,
	MaxPathDepth:           2,
	AllowedExternalDomains: []string{"cdn.example.com"},
	SensitivePatterns:      []string{"admin", "api", "private"},
	AllowedExtensions:      []string{".txt", ".xml", ".json", ""},
	Workers:                3,
}

const (
	robotsDomain           = "www.example.com"
	robotsRegisteredDomain = "example.com"
)

// TestExtractUrlsFromRobots tests URL extraction from robots.txt content
func TestExtractUrlsFromRobots(t *testing.T) {
	ctx := context.Background()
	q := queue.NewQueue(1)

	robotsContent := `
# Robots.txt example
User-agent: *

Disallow: /docs/guide
Allow: /public/data
Disallow: /search?q=
Disallow: /api/v1/users
Disallow: https://example.com/api/v1
Disallow: /path1/path2/path3
Sitemap: //cdn.example.com/sitemap.xml
`

	extractUrlsFromRobots(ctx, robotsContent, robotsTestConfig, q, robotsDomain, robotsRegisteredDomain)
	time.Sleep(10 * time.Millisecond)

	extractedUrls := q.GetFoundUrls()
	slices.Sort(extractedUrls)

	// Expected URLs after validation
	// Note: robotsDomain is www.example.com, so relative URLs get that domain
	// Sitemap URL not extracted - appears to be a parser issue
	expectedUrls := []string{
		"https://example.com/api/v1",
		"https://cdn.example.com/sitemap.xml",
		"https://www.example.com/api/v1/users",
		"https://www.example.com/docs/guide",
		"https://www.example.com/public/data",
		"https://www.example.com/search?q=",
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
