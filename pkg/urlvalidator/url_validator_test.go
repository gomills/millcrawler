package urlvalidator

import (
	"testing"

	"github.com/gomills/gofocusedcrawler/internal/config"
)

var (
	testConfig = &config.Config{
		TimeOutDuration:        60,
		CustomUserAgent:        "Test",
		MaxPathDepth:           1,
		AllowedExternalDomains: []string{"github.com"},
		SensitivePatterns:      []string{"dashboard", "test", "repo", "private"},
		AllowedExtensions:      []string{".git", ".txt", ".js", ".html", ""},
		Workers:                3,
	}
	domain           = "www.example.com"
	registeredDomain = "example.com"
)

var validAbsoluteUrls = []string{
	"https://example.com/a",
	"https://example.com/bb/a/3f/sf.txt",
	"//example.com/a/b/c.js",
	"//example.com/a/b/fsf/c.git",
	"example.com/a",
	"https://example.com/fe/s/afe/safesa/dashboard",
	"https://example.com/test",
	"https://github.com/private/repo/example",
	"https://example.com/a/b/dashboard",
	"https://testing.example.com/a/b/dashboard",
}

var invalidAbsoluteUrls = []string{
	"https://example.com/a/b/",
	"https://cdn.example.com/libs/jquery.min.js",
	"https://example.com/jquery.min.js",
	"https://example.com/a/b",
	"https://example.com/a/b/c#section",
	"https://youtube.com/a/b/c",
	"https://youtube.com",
	"//example.com/a/b/c.html",
	"example.com/a/b/c.htm",
	"fes ojfeio oiejfisofjesio",
	"https is a internet protocol",
	"https://www.youtube.com/embed/ZipKoVUSWlY?rel=0&modestbranding=1&wmode=opaque",
}

var validRelativeUrls = []string{
	"/a/b/c.js?query=1111",
	"/a/b/c.txt",
	"/afesfesf",
	"/?q=user/login",
	"/esfes/private?query",
}

var invalidRelativeUrls = []string{
	"/a/b/c",
	"/a/b/c*extra",
	"/fesfs/fesfe/fesfefs/v",
	"/ fesioj ifoesjifoe",
}

func TestValidateStringForUrl(t *testing.T) {
	// Test valid absolute URLs
	for _, u := range validAbsoluteUrls {
		url, err := ValidateStringForUrl(testConfig, u, domain, registeredDomain)
		if err != nil || url == nil {
			t.Errorf("Expected valid URL %s, got error: %v", u, err)
		}
	}

	// Test invalid absolute URLs
	for _, u := range invalidAbsoluteUrls {
		_, err := ValidateStringForUrl(testConfig, u, domain, registeredDomain)
		if err == nil {
			t.Errorf("Expected invalid URL %s to fail validation", u)
		}
	}

	// Test valid relative URLs
	for _, u := range validRelativeUrls {
		url, err := ValidateStringForUrl(testConfig, u, domain, registeredDomain)
		if err != nil || url == nil {
			t.Errorf("Expected valid relative URL %s, got error: %v", u, err)
		}
	}

	// Test invalid relative URLs
	for _, u := range invalidRelativeUrls {
		_, err := ValidateStringForUrl(testConfig, u, domain, registeredDomain)
		if err == nil {
			t.Errorf("Expected invalid relative URL %s to fail validation", u)
		}
	}
}
