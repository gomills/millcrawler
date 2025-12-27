package crawler

import (
	"fmt"

	http "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
	"github.com/bogdanfinn/tls-client/profiles"
)

// getStealthHttpsClient returns an HTTP client with custom TLS fingerprint and according headers
// https://bogdanfinn.gitbook.io/open-source-oasis/tls-client/installation-and-quick-usage
func getStealthHttps() (tls_client.HttpClient, http.Header, error) {

	// generate client
	options := []tls_client.HttpClientOption{
		tls_client.WithTimeoutSeconds(15),
		tls_client.WithClientProfile(profiles.Firefox_123),
	}

	client, err := tls_client.NewHttpClient(tls_client.NewNoopLogger(), options...)
	if err != nil {
		return nil, nil, fmt.Errorf("client_gen_err")
	}

	// generate according headers. IMPORTANT: match browser headers
	headers := http.Header{
		"accept":          {"*/*"},
		"accept-language": {"en-US;q=0.8,en;q=0.7"},
		"user-agent":      {"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.75 Safari/537.36"},
		http.HeaderOrderKey: {
			"accept",
			"accept-language",
			"user-agent",
		},
	}

	return client, headers, nil
}
