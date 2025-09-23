package urlvalidator

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/gomills/gofocusedcrawler/internal/config"
	"golang.org/x/net/publicsuffix"
)

// ValidateStringForUrl validates a string as a URL based on heuristics, cleans and parses it
func ValidateStringForUrl(config *config.Config, possibleUrl string, domain string, registeredDomain string) (*url.URL, error) {

	possibleUrl = strings.TrimSpace(possibleUrl)

	// Pass heuristic check
	err := initialCheck(possibleUrl)
	if err != nil {
		return nil, err
	}

	// Resolve it, add protocol, clean in general and parse it
	cleanParsedUrl, isLocal, err := cleanAndParseUrl(possibleUrl, domain, registeredDomain)
	if err != nil {
		return nil, err
	}

	// Decide if it's of use for us based on heuristics
	if isLocal {
		return ValidateLocalUrl(config, cleanParsedUrl, domain, registeredDomain)
	} else {
		return validateExternalUrl(config, cleanParsedUrl)
	}

}

// initialCheck performs basic heuristic checks on the possible URL string.
func initialCheck(possibleUrl string) error {
	if possibleUrl == "" || len([]rune(possibleUrl)) < 3 || len(possibleUrl) > 300 {
		return fmt.Errorf("Didn't pass length check")
	} else if strings.Contains(possibleUrl, " ") || strings.Contains(possibleUrl, "\n") || strings.HasPrefix(possibleUrl, "mailto") {
		return fmt.Errorf("Possibly not a URL")
	} else {
		return nil
	}

}

// cleanAndParseUrl cleans the URL string and parses it, returning the parsed URL and locality.
func cleanAndParseUrl(possibleUrl string, domain string, registeredDomain string) (*url.URL, bool, error) {
	cleanedUrl := cleanUrl(possibleUrl)
	return solveUrl(cleanedUrl, domain, registeredDomain)
}

// cleanUrl removes fragments, wildcards and adds protocol in case of a protocol-relative url (e.g.: "//example.com/path1")
func cleanUrl(possibleUrl string) string {
	if hsh := strings.IndexRune(possibleUrl, '#'); hsh != -1 {
		possibleUrl = possibleUrl[:hsh]
	}
	if strk := strings.IndexRune(possibleUrl, '*'); strk != -1 {
		possibleUrl = possibleUrl[:strk]
	}
	if strings.HasPrefix(possibleUrl, `//`) {
		possibleUrl = "https:" + possibleUrl
	}

	return possibleUrl
}

// solveUrl parses the cleaned URL, resolves relative URLs, and determines if it is local.
// Returns (parsed url, isLocal bool, error).
func solveUrl(cleanedUrl string, domain string, registeredDomain string) (*url.URL, bool, error) {

	// Solve relative URLs
	if strings.HasPrefix(cleanedUrl, `/`) {

		cleanedUrl = `https://` + domain + cleanedUrl

		cleanedParsedUrl, err := url.Parse(cleanedUrl)
		if err != nil {
			return nil, true, err
		}

		return cleanedParsedUrl, true, nil
	}

	// Add protocol if missing in cases like example.com/path1/path2. Protocol-relative URLs are handled in cleanUrl()
	if !strings.HasPrefix(cleanedUrl, "http") {
		cleanedUrl = `https://` + cleanedUrl
	}

	cleanedParsedUrl, pErr := url.Parse(cleanedUrl)
	if cleanedParsedUrl == nil || pErr != nil {
		return nil, false, pErr
	}

	foundUrlRegDomain, psErr := publicsuffix.EffectiveTLDPlusOne(cleanedParsedUrl.Hostname())
	if psErr != nil {
		return nil, false, psErr
	}

	if foundUrlRegDomain == registeredDomain {
		return cleanedParsedUrl, true, nil
	} else {
		return cleanedParsedUrl, false, nil
	}

}
