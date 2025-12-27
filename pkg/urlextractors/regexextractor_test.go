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

	// Text with 10 URLs mixed in random text
	// 6 valid: 3 absolute, 3 relative
	// 4 invalid: missing extensions or sensitive patterns
	textContent := `
Lorem ipsum dolor sit amet, consectetur adipiscing elit. 
Visit //example.com/products for our latest offerings.
Sed do eiusmod tempor incididunt ut labore https://cdn.example.com/assets/style.css et dolore magna aliqua.
Ut enim ad minim veniam, quis nostrud exercitation ullamco https://github.com/projects/repo laboris nisi ut aliquip ex ea commodo consequat.

Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur.
Check https://example.com/api/v1/users for user management.
Excepteur sint occaecat cupidatat non proident, sunt in culpa //example.com/docs/guide qui officia deserunt mollit anim id est laborum.

Some invalid URLs that should not be extracted:
- https://example.com/dasfseard.esfpso 
- https://example.com/admfein/seettings.fef 
- //example.com/privaesefte/dfeata.css 
- //example.com/protecfested.exe

Additional valid URLs scattered throughout:
Check our API at https://cdn.example.com/docs and support at //github.com/support/issues.
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
		"https://example.com/products",
		"https://github.com/projects/repo",
		"https://example.com/api/v1/users",
		"https://example.com/docs/guide qui",
		"https://cdn.example.com/docs",
		"https://github.com/support/issues",
	}
	slices.Sort(expectedUrls)

	if len(extractedUrls) != len(expectedUrls) {
		t.Logf("Extracted URLs (%d): %v", len(extractedUrls), extractedUrls)
		t.Logf("Expected URLs (%d): %v", len(expectedUrls), expectedUrls)
	}

	for i, expected := range expectedUrls {
		if i >= len(extractedUrls) || (extractedUrls[i] != expected) {
			t.Errorf("URL mismatch at index %d: expected %q, got %q", i, expected, extractedUrls[i])
		}
	}
}
