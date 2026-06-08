package crawler

import (
	http "github.com/bogdanfinn/fhttp"
	"github.com/gomills/millcrawler/pkg/queue"
	"github.com/gomills/millcrawler/pkg/secrets"
)

// scanSecrets takes a response's body and headers and scans them for secrets using the secrets package.
// Any found secret is added to the queue's secret vault.
func scanSecrets(qp *queue.Queue, headers http.Header, url string, bodyBytes []byte) {
	// scan headers for secrets
	for _, values := range headers {
		for _, value := range values {

			matches, _ := secrets.FindSecrets(value)
			for _, match := range matches {
				qp.AddSecret(match, url)
			}

		}
	}

	// scan body for secrets
	bodyStr := string(bodyBytes)
	matches, _ := secrets.FindSecrets(bodyStr)
	for _, match := range matches {
		qp.AddSecret(match, url)
	}
}
