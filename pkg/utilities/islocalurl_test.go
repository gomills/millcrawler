package utilities

import (
	"net/url"
	"testing"
)

func TestIsLocalUrl(t *testing.T) {
	tests := []struct {
		name             string
		url              string
		registeredDomain string
		want             bool
		wantErr          bool
	}{
		{
			name:             "same domain",
			url:              "https://example.com/path",
			registeredDomain: "example.com",
			want:             true,
			wantErr:          false,
		},
		{
			name:             "subdomain of same domain",
			url:              "https://api.example.com/path",
			registeredDomain: "example.com",
			want:             true,
			wantErr:          false,
		},
		{
			name:             "different domain",
			url:              "https://google.com/path",
			registeredDomain: "example.com",
			want:             false,
			wantErr:          false,
		},
		{
			name:             "subdomain of different domain",
			url:              "https://api.google.com/path",
			registeredDomain: "example.com",
			want:             false,
			wantErr:          false,
		},
		{
			name:             "localhost invalid for registered domain",
			url:              "http://localhost/path",
			registeredDomain: "example.com",
			want:             false,
			wantErr:          true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsedUrl, err := url.Parse(tt.url)
			if err != nil {
				t.Fatalf("Failed to parse test URL: %v", err)
			}

			got, err := IsLocalUrl(parsedUrl, tt.registeredDomain)

			if (err != nil) != tt.wantErr {
				t.Errorf("IsLocalUrl() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("IsLocalUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}
