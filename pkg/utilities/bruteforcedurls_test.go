package utilities

import (
	"net/url"
	"slices"
	"testing"

	"github.com/gomills/millcrawler/pkg/config"
)

func TestGetBruteForcedUrls(t *testing.T) {
	tests := []struct {
		name             string
		domain           string
		registeredDomain string
		wantErr          bool
		wantURLCount     int
		wantURLsContain  []string
	}{
		{
			name:             "valid domains",
			domain:           "example.com",
			registeredDomain: "example.com",
			wantErr:          false,
			wantURLCount:     8,
			wantURLsContain: []string{
				"https://example.com/",
				"https://example.com/robots.txt",
				"https://example.com/sitemap.xml",
				"https://dev.example.com/",
				"https://staging.example.com/",
				"https://admin.example.com/",
				"https://test.example.com/",
				"https://internal.example.com/",
			},
		},
		{
			name:             "subdomain as domain",
			domain:           "api.example.com",
			registeredDomain: "example.com",
			wantErr:          false,
			wantURLCount:     8,
			wantURLsContain: []string{
				"https://api.example.com/",
				"https://dev.example.com/",
			},
		},
		{
			name:             "empty domain",
			domain:           "",
			registeredDomain: "example.com",
			wantErr:          true,
		},
		{
			name:             "empty registered domain",
			domain:           "example.com",
			registeredDomain: "",
			wantErr:          true,
		},
		{
			name:             "both empty",
			domain:           "",
			registeredDomain: "",
			wantErr:          true,
		},
		{
			name:             "whitespace only domain",
			domain:           "   ",
			registeredDomain: "example.com",
			wantErr:          true,
		},
		{
			name:             "whitespace only registered domain",
			domain:           "example.com",
			registeredDomain: "   ",
			wantErr:          true,
		},
	}

	config, _ := config.LoadConfig(10, 1, 1, []string{}, []string{}, []string{}, 10, true, true, true, true)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetBruteForcedUrls(tt.domain, tt.registeredDomain, config)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetBruteForcedUrls() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			if len(got) != tt.wantURLCount {
				t.Errorf("GetBruteForcedUrls() returned %d URLs, want %d", len(got), tt.wantURLCount)
			}

			gotURLStrs := make([]string, len(got))
			for i, u := range got {
				gotURLStrs[i] = u.String()
			}

			for _, want := range tt.wantURLsContain {
				found := false
				for _, gotStr := range gotURLStrs {
					if gotStr == want {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("GetBruteForcedUrls() missing expected URL: %s\ngot: %v", want, gotURLStrs)
				}
			}
		})
	}
}

func TestGetBruteForcedUrlsURLValidity(t *testing.T) {
	// Test that all returned URLs are valid and properly formed
	config, _ := config.LoadConfig(10, 1, 1, []string{}, []string{}, []string{}, 10, true, true, true, true)
	got, err := GetBruteForcedUrls("example.com", "example.com", config)
	if err != nil {
		t.Fatalf("GetBruteForcedUrls() unexpected error: %v", err)
	}

	for i, u := range got {
		// Check that URL is not nil
		if u == nil {
			t.Errorf("URL at index %d is nil", i)
			continue
		}

		// Check that scheme is https
		if u.Scheme != "https" {
			t.Errorf("URL at index %d has scheme %q, want https", i, u.Scheme)
		}

		// Check that host is set
		if u.Host == "" {
			t.Errorf("URL at index %d has empty host", i)
		}

		// Ensure URL can be parsed again (round-trip test)
		reparsed, err := url.Parse(u.String())
		if err != nil {
			t.Errorf("URL at index %d failed to round-trip parse: %v", i, err)
		}

		if reparsed.String() != u.String() {
			t.Errorf("URL at index %d round-trip mismatch:\n  original: %s\n  reparsed: %s", i, u.String(), reparsed.String())
		}
	}
}

func TestGetBruteForcedUrlsNoDuplicates(t *testing.T) {
	config, _ := config.LoadConfig(10, 1, 1, []string{}, []string{}, []string{}, 10, true, true, true, true)
	got, err := GetBruteForcedUrls("example.com", "example.com", config)
	if err != nil {
		t.Fatalf("GetBruteForcedUrls() unexpected error: %v", err)
	}

	seen := make(map[string]bool)
	for i, u := range got {
		uStr := u.String()
		if seen[uStr] {
			t.Errorf("Duplicate URL found at index %d: %s", i, uStr)
		}
		seen[uStr] = true
	}
}

func TestSubdomainsBruteforceBoolFlagWorks(t *testing.T) {
	domain := "example.com"
	registeredDomain := domain

	bruteForcedPaths := []string{
		"https://example.com/",
		"https://example.com/robots.txt",
		"https://example.com/sitemap.xml",
	}
	bruteForcedSubdomains := []string{
		"https://dev.example.com/",
		"https://staging.example.com/",
		"https://admin.example.com/",
		"https://test.example.com/",
		"https://internal.example.com/",
	}

	// this should make them
	configurationWithSubdomains, _ := config.LoadConfig(10, 1, 1, []string{}, []string{}, []string{}, 10, true, true, true, true)
	seedUrlsWithSubdomains, _ := GetBruteForcedUrls(domain, registeredDomain, configurationWithSubdomains)

	// here no subdomains should be bruteforced
	configurationWithoutSubdomains, _ := config.LoadConfig(10, 1, 1, []string{}, []string{}, []string{}, 10, true, false, true, true)
	seedUrlsWithoutSubdomains, _ := GetBruteForcedUrls(domain, registeredDomain, configurationWithoutSubdomains)

	for _, u := range bruteForcedPaths {

		if !slices.ContainsFunc(seedUrlsWithSubdomains, func(x *url.URL) bool { return x.String() == u }) {
			t.Errorf("Bruteforced urls doesn't contain at least url %s", u)
		}
		if !slices.ContainsFunc(seedUrlsWithoutSubdomains, func(x *url.URL) bool { return x.String() == u }) {
			t.Errorf("Bruteforced urls doesn't contain at least url %s", u)
		}
	}

	for _, u := range bruteForcedSubdomains {

		if !slices.ContainsFunc(seedUrlsWithSubdomains, func(x *url.URL) bool { return x.String() == u }) {
			t.Errorf("Bruteforced urls with subdomains doesn't contain at least subdomain %s", u)
		}
		if slices.ContainsFunc(seedUrlsWithoutSubdomains, func(x *url.URL) bool { return x.String() == u }) {
			t.Errorf("Bruteforced urls without subdomains contains at least subdomain %s", u)
		}
	}

}
