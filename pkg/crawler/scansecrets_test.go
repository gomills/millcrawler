package crawler

import (
	"testing"

	http "github.com/bogdanfinn/fhttp"
	"github.com/gomills/millcrawler/pkg/queue"
)

// TestScanSecrets_Headers_NoSecrets verifies that scanSecrets correctly handles headers with no secrets
func TestScanSecrets_Headers_NoSecrets(t *testing.T) {
	q := queue.NewQueue(1)
	url := "https://example.com/test"

	headers := http.Header{
		"Content-Type": []string{"application/json"},
		"User-Agent":   []string{"Mozilla/5.0"},
	}

	bodyBytes := []byte("")

	scanSecrets(q, headers, url, bodyBytes)

	if count := q.GetNSecrets(); count != 0 {
		t.Errorf("Expected 0 secrets, got %d", count)
	}
}

// TestScanSecrets_Headers_OneSecret verifies that scanSecrets correctly extracts one secret from headers
func TestScanSecrets_Headers_OneSecret(t *testing.T) {
	q := queue.NewQueue(1)
	url := "https://example.com/test"

	headers := http.Header{
		"Content-Type":  []string{"application/json"},
		"Authorization": []string{"Bearer mywebsecret_VzQiRPwmcKLVu72PMizPpuYc"},
	}

	bodyBytes := []byte("")

	scanSecrets(q, headers, url, bodyBytes)

	if count := q.GetNSecrets(); count != 1 {
		t.Errorf("Expected 1 secret, got %d", count)
	}

	secrets := q.GetSecretsCopy()
	if storedUrl, ok := secrets["mywebsecret_VzQiRPwmcKLVu72PMizPpuYc"]; !ok {
		t.Errorf("Expected secret mywebsecret_VzQiRPwmcKLVu72PMizPpuYc to be stored")
	} else if storedUrl != url {
		t.Errorf("Expected URL %q, got %q", url, storedUrl)
	}
}

// TestScanSecrets_Headers_TwoSecrets verifies that scanSecrets correctly extracts two secrets from headers
func TestScanSecrets_Headers_TwoSecrets(t *testing.T) {
	q := queue.NewQueue(1)
	url := "https://example.com/test"

	headers := http.Header{
		"X-Api-Key":     []string{"mywebsecret_XyZaBcDeFgHiJkLmNoPqRs"},
		"Authorization": []string{"Bearer mywebsecret_VzQiRPwmcKLVu72PMizPpuYc"},
	}

	bodyBytes := []byte("")

	scanSecrets(q, headers, url, bodyBytes)

	if count := q.GetNSecrets(); count != 2 {
		t.Errorf("Expected 2 secrets, got %d", count)
	}

	secrets := q.GetSecretsCopy()

	expectedSecrets := map[string]string{
		"mywebsecret_XyZaBcDeFgHiJkLmNoPqRs":   url,
		"mywebsecret_VzQiRPwmcKLVu72PMizPpuYc": url,
	}

	for secret, expectedUrl := range expectedSecrets {
		if storedUrl, ok := secrets[secret]; !ok {
			t.Errorf("Expected secret %q to be stored", secret)
		} else if storedUrl != expectedUrl {
			t.Errorf("Expected URL %q for secret %q, got %q", expectedUrl, secret, storedUrl)
		}
	}
}

// TestScanSecrets_Body_NoSecrets verifies that scanSecrets correctly handles body with no secrets
func TestScanSecrets_Body_NoSecrets(t *testing.T) {
	q := queue.NewQueue(1)
	url := "https://example.com/test"

	headers := http.Header{}
	bodyBytes := []byte("This is a normal response body with no secrets")

	scanSecrets(q, headers, url, bodyBytes)

	if count := q.GetNSecrets(); count != 0 {
		t.Errorf("Expected 0 secrets, got %d", count)
	}
}

// TestScanSecrets_Body_OneSecret verifies that scanSecrets correctly extracts one secret from body
func TestScanSecrets_Body_OneSecret(t *testing.T) {
	q := queue.NewQueue(1)
	url := "https://example.com/test"

	headers := http.Header{}
	bodyBytes := []byte(`
Configuration:
api_key: mywebsecret_AbCdEfGhIjKlMnOpQrStUv
database: postgres
`)

	scanSecrets(q, headers, url, bodyBytes)

	if count := q.GetNSecrets(); count != 1 {
		t.Errorf("Expected 1 secret, got %d", count)
	}

	secrets := q.GetSecretsCopy()
	if storedUrl, ok := secrets["mywebsecret_AbCdEfGhIjKlMnOpQrStUv"]; !ok {
		t.Errorf("Expected secret mywebsecret_AbCdEfGhIjKlMnOpQrStUv to be stored")
	} else if storedUrl != url {
		t.Errorf("Expected URL %q, got %q", url, storedUrl)
	}
}

// TestScanSecrets_Body_TwoSecrets verifies that scanSecrets correctly extracts two secrets from body
func TestScanSecrets_Body_TwoSecrets(t *testing.T) {
	q := queue.NewQueue(1)
	url := "https://example.com/test"

	headers := http.Header{}
	bodyBytes := []byte(`
Credentials:
mysuperweb_key: mywebsecret_KeyAbCdEfGhIjKlMn
api_token: mywebsecret_ApiTokenXyZaBcDeFgHiJkLmNo
database: postgres
`)

	scanSecrets(q, headers, url, bodyBytes)

	if count := q.GetNSecrets(); count != 2 {
		t.Errorf("Expected 2 secrets, got %d", count)
	}

	secrets := q.GetSecretsCopy()

	expectedSecrets := map[string]string{
		"mywebsecret_KeyAbCdEfGhIjKlMn":          url,
		"mywebsecret_ApiTokenXyZaBcDeFgHiJkLmNo": url,
	}

	for secret, expectedUrl := range expectedSecrets {
		if storedUrl, ok := secrets[secret]; !ok {
			t.Errorf("Expected secret %q to be stored", secret)
		} else if storedUrl != expectedUrl {
			t.Errorf("Expected URL %q for secret %q, got %q", expectedUrl, secret, storedUrl)
		}
	}
}
