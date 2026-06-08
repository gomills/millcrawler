package secrets

import (
	"regexp"
	"strings"
)

// secretStruct contains the regex to detect the secret and a prefix pattern for a prior substring search.
// Prefix should be zero-value-string if it doesn't exist.
type secretStruct struct {
	prefix string
	regex  *regexp.Regexp
}

// here we declare the actual secrets structs
var (
	sampleSecret = secretStruct{prefix: "mywebsecret_", regex: regexp.MustCompile(`mywebsecret_[a-zA-Z0-9]{16,247}`)}
)

// here we collect them all together for iteration
var (
	allSecretsStructs = []secretStruct{sampleSecret}
)

// FindSecrets iterates over the hardcoded allSecretsStructs to match potential secrets in the given content.
// It first performs a substring search on prefixed patterns for efficiency, then applies the regex.
// For non-prefixed patterns, it goes directly to regex matching.
func FindSecrets(content string) ([]string, error) {

	var allMatches []string

	// iterate through all secrets structs collection
	for _, sp := range allSecretsStructs {

		// if the pattern has a prefix, proceed to substring check first
		if sp.prefix != "" {

			if !strings.Contains(content, sp.prefix) {
				continue
			}

		}

		// if prefix found or no prefix, apply the regex
		matches := sp.regex.FindAllString(content, -1)
		allMatches = append(allMatches, matches...)
	}

	return allMatches, nil
}
