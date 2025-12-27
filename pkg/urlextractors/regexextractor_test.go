package urlsextractors

import (
	"context"
	"io"
	"slices"
	"strings"
	"testing"
	"time"

	http "github.com/bogdanfinn/fhttp"
	"github.com/gomills/gofocusedcrawler/internal/config"
	"github.com/gomills/gofocusedcrawler/pkg/queue"
)

var regexTestConfig = &config.Config{
	TimeOutDuration:        60,
	MaxPathDepth:           2,
	AllowedExternalDomains: []string{"github.com", "cdn.example.com"},
	SensitivePatterns:      []string{"dashboard", "admin", "private"},
	AllowedExtensions:      []string{".html", ".txt", ".json", ""},
	Workers:                3,
}

const (
	regexDomain           = "www.example.com"
	regexRegisteredDomain = "example.com"
)

// TestRegexFileForUrls tests regex extraction with mixed valid and invalid URLs embedded in text
func TestRegexFileForUrls(t *testing.T) {
	ctx := context.Background()
	q := queue.NewQueue(1)

	textContent := `
Lorem ipsum dolor sit amet, consectetur adipiscing elit. 
Visit https://example.com/path1 for our latest offerings. 
Sed do eiusmod tempor incididunt ut labore 
https://cdn.example.com/path1/path2/path3.json et dolore magna aliqua.
Ut enim ad minim veniam, quis nostrud exercitation ullamco 
https://github.com/path1/path2 laboris nisi ut 
aliquip ex ea commodo consequat.

Duis aute irure dolor in reprehenderit in voluptate velit 
esse cillum dolore eu fugiat nulla pariatur.
Check https://example.com/path1/path2/path3 for user management.
Check https://example.com/path1/path2/path3/path4 for user management.
Excepteur sint occaecat cupidatat non proident, sunt in culpa 
https://example.com/docs/guide qui officia deserunt mollit anim id est laborum.

Qui officia deserunt mollit anim id est laborum.:
- https://example.com/dashboard/a/a.json 
- https://example.com/admin/settings.txt 
- //example.com/privaesefte/dfeata.css

https://example.com/path1/admin/path2/caca)

https://example.com/path1/private/path2/culo'

Excepteur sint occaecat cupidatat:
Check our API at https://cdn.example.com/docs 
and support at //github.com/support/issues.
`

	response := &http.Response{
		Body: io.NopCloser(strings.NewReader(textContent)),
	}

	regexFileForUrls(ctx, regexTestConfig, response, q, regexDomain, regexRegisteredDomain)

	time.Sleep(50 * time.Millisecond)

	extractedUrls := q.GetFoundUrls()
	slices.Sort(extractedUrls)

	// Expected valid URLs (after validation and deduplication)
	expectedUrls := []string{
		"https://example.com/path1",
		"https://cdn.example.com/path1/path2/path3.json",
		"https://github.com/path1/path2",
		"https://example.com/docs/guide",
		"https://example.com/dashboard/a/a.json",
		"https://example.com/admin/settings.txt",
		"https://example.com/path1/admin/path2/caca",
		"https://example.com/path1/private/path2/culo",
		"https://cdn.example.com/docs",
		"https://github.com/support/issues",
	}
	slices.Sort(expectedUrls)

	if len(extractedUrls) != len(expectedUrls) {
		t.Logf("Extracted URLs (%d)", len(extractedUrls))
		t.Logf("Expected URLs (%d)", len(expectedUrls))
	}

	for i, expected := range expectedUrls {
		if i >= len(extractedUrls) || (extractedUrls[i] != expected) {
			t.Errorf("URL mismatch at index %d: expected %q, got %q", i, expected, extractedUrls[i])
		}
	}
}
