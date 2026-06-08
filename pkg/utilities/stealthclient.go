package utilities

import (
	http "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
	"github.com/bogdanfinn/tls-client/profiles"
	"github.com/gomills/millcrawler/pkg/config"
)

const (
	chromeUserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.75 Safari/537.36"
)

// getStealthHttpsClient returns an HTTP client with custom TLS fingerprint and according headers
// Example from source library: https://bogdanfinn.gitbook.io/open-source-oasis/tls-client/installation-and-quick-usage
// available and valid clients: https://bogdanfinn.gitbook.io/open-source-oasis/tls-client/supported-and-tested-client-profiles
func GetStealthHttps(config *config.Config) (tls_client.HttpClient, error) {

	var profile profiles.ClientProfile
	var userAgent string

	profile = profiles.Chrome_133
	userAgent = chromeUserAgent

	// generate according headers. It's in the TODO.md list to further improve these headers
	headers := http.Header{
		"accept":          {"*/*"},
		"accept-language": {"en-US,en;q=0.5"},
		"user-agent":      {userAgent},
		http.HeaderOrderKey: {
			"accept",
			"accept-language",
			"user-agent",
		},
	}

	// craft options for client
	clientOptions := []tls_client.HttpClientOption{
		tls_client.WithTimeoutSeconds(15),
		tls_client.WithClientProfile(profile),
		tls_client.WithDefaultHeaders(headers),
	}

	// add a cookies jar if config says so
	if config.Cookies {
		clientOptions = append(clientOptions, tls_client.WithCookieJar(tls_client.NewCookieJar()))
	}

	client, err := tls_client.NewHttpClient(tls_client.NewNoopLogger(), clientOptions...)
	if err != nil {
		return nil, err
	}

	return client, nil
}
