package crawler

import (
	"fmt"
	"math/rand"

	http "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
	"github.com/bogdanfinn/tls-client/profiles"
)

const (
	firefoxUsAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:120.0) Gecko/20100101 Firefox/120.0"
	chromeUsAgent  = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.75 Safari/537.36"
)

// getStealthHttpsClient returns an HTTP client with custom TLS fingerprint and according headers
// https://bogdanfinn.gitbook.io/open-source-oasis/tls-client/installation-and-quick-usage
func getStealthHttps() (tls_client.HttpClient, http.Header, error) {

	// randomly choose between Chrome and Firefox
	useChrome := rand.Intn(2) == 0

	var profile profiles.ClientProfile
	var userAgent string

	if useChrome {
		profile = profiles.Chrome_120
		userAgent = chromeUsAgent
	} else {
		profile = profiles.Firefox_120
		userAgent = firefoxUsAgent
	}

	// generate client
	options := []tls_client.HttpClientOption{
		tls_client.WithTimeoutSeconds(15),
		tls_client.WithClientProfile(profile),
	}

	client, err := tls_client.NewHttpClient(tls_client.NewNoopLogger(), options...)
	if err != nil {
		return nil, nil, fmt.Errorf("client_gen_err")
	}

	// generate according headers. IMPORTANT: match browser headers
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

	return client, headers, nil
}
