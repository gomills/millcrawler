package urlvalidator

import (
	"testing"

	"github.com/gomills/gofocusedcrawler/internal/config"
)

var (
	testConfig = &config.Config{
		TimeOutDuration:        60,
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
	// Basic local URLs with allowed extensions
	"https://example.com/a",
	"https://example.com/bb/a/3f/sf.txt",
	"//example.com/a/b/c.js",
	"//example.com/a/b/fsf/c.git",
	"example.com/a",

	// Local URLs with sensitive patterns (bypass depth check)
	"https://example.com/fe/s/afe/safesa/dashboard",
	"https://example.com/test",
	"https://example.com/a/b/dashboard",
	"https://testing.example.com/a/b/dashboard",

	// External URLs in allowed list
	"https://github.com/private/repo/example",

	// Root paths (depth=0)
	"https://example.com/",
	"https://example.com",
	"example.com/",
	"//example.com/",

	// Single component paths (depth=1, at MaxPathDepth limit)
	"https://example.com/path",
	"https://example.com/api.txt",
	"https://example.com/data.git",

	// Subdomains (should be treated as local)
	"https://api.example.com/endpoint",
	"https://cdn.example.com/data.txt",
	"https://dev.example.com/test",
	"https://staging.example.com/normal",

	// Query strings with allowed extensions
	"https://example.com/a?query=value",
	"https://example.com/file.js?v=1.0",

	// Fragment removal (should work)
	"https://example.com/a#section",
	"https://example.com/test#anchor",
	"https://example.com/test*",

	// Multiple sensitive patterns in deep paths
	"https://example.com/a/b/c/d/private/test",
	"https://example.com/repo/dashboard/aha/huhu",

	// Protocol-relative URLs
	"//testing.example.com/data",
	"//api.example.com/config.txt",
}

var invalidAbsoluteUrls = []string{
	// No allowed extension (exceeds depth without sensitive pattern)
	"https://example.com/a/b/",
	"https://example.com/a/b",
	"https://example.com/a/b/c#section",

	// Common JS libraries (should be discarded)
	"https://cdn.example.com/libs/jquery.min.js",
	"https://example.com/jquery.min.js",
	"https://example.com/vendor/bootstrap.min.js",
	"https://example.com/node_modules/react.js",
	"https://example.com/libs/angular.js",

	// Disallowed extensions
	"https://example.com/a.png",
	"https://example.com/image.svg",
	"https://example.com/style.css",
	"https://example.com/data.json",
	"https://example.com/video.mp4",

	// External domains not in allowed list
	"https://youtube.com/a/b/c",
	"https://youtube.com",
	"https://www.google.com/search",
	"https://facebook.com/user/profile",

	// HTML files exceeding depth without sensitive pattern
	"//example.com/a/b/c.html",
	"example.com/a/b/c.htm",
	"https://example.com/a/b/c.html",

	// Invalid characters/patterns
	"fes ojfeio oiejfisofjesio",
	"https is a internet protocol",
	"https://www.youtube.com/embed/ZipKoVUSWlY?rel=0&modestbranding=1&wmode=opaque",

	// Too short (less than 3 chars)
	"ab",
	"//",
	"",

	// Contains whitespace or newlines
	"https://example.com/a b",
	"https://example.com/path\ninjection",

	// Mailto links
	"mailto:user@example.com",
	"mailto://example.com",

	// Empty hostname (protocol-only)
	"https://",
	"http://",

	// Unrecognized schemes (treated as no scheme and prefixed with https://)
	"ftp://files.example.com/data",
	"ws://example.com/socket",

	// Too large, > 600 characters
	"https://example.com/aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
}

var validRelativeUrls = []string{
	// Root path (with content to satisfy 3-char minimum)
	"/?x=1",

	// Single component (depth=1, at MaxPathDepth limit)
	"/abc",
	"/index.html",
	"/data.txt",
	"/config.js",

	// With allowed extensions
	"/a/b/c.js?query=1111",
	"/a/b/c.txt",

	// With query parameters
	"/?q=user/login",
	"/path?param=value",

	// With sensitive patterns (bypass depth check)
	"/afesfesf",
	"/esfes/private?query",
	"/deep/path/with/dashboard/and/test",
	"/repo/private/dashboard",
	"/admin/test/repo/user",

	// Root with trailing slash and query
	"/?query=test",

	// Fragment handling
	"/page#section",
	"/doc.html#anchor",
}

var invalidRelativeUrls = []string{
	// Exceeds depth without sensitive pattern
	"/a/b/c",
	"/fesfs/fesfe/fesfefs/v",
	"/deep/path/without/patterns",

	// Wildcard in path (truncated, may become valid - skip uncertain cases)

	// Contains whitespace
	"/ fesioj ifoesjifoe",
	"/path with spaces",
	"/bad\npath",

	// Common JS libraries
	"/vendor/jquery.min.js",
	"/node_modules/react.js",
	"/libs/bootstrap.js",
	"/js/angular.js",

	// Disallowed extensions
	"/image.png",
	"/style.css",
	"/data.json",
	"/video.mp4",
	"/archive.zip",

	// Too short
	"a",
	"",
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
