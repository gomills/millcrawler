package urlvalidator

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/gomills/gofocusedcrawler/internal/config"
	"golang.org/x/net/publicsuffix"
)

// ValidateStringForUrl validates a URL as a string based on heuristics, cleans and parses it.
// Returns not-nil error if not valid.
func ValidateStringForUrl(config *config.Config, possibleUrl string, domain string, registeredDomain string) (*url.URL, error) {

	// 0. trim whitespaces
	possibleUrl = strings.TrimSpace(possibleUrl)

	// 1. pass heuristic check
	err := initialCheck(possibleUrl)
	if err != nil {
		return nil, err
	}

	// 2. resolve it, add protocol, clean in general and parse it
	cleanParsedUrl, isLocal, err := cleanAndParseUrl(possibleUrl, domain, registeredDomain)
	if err != nil {
		return nil, err
	}

	// 3. decide if it passes heuristics to be of use
	if isLocal {
		return validateLocalUrl(config, cleanParsedUrl, domain, registeredDomain)
	} else {
		return validateExternalUrl(config, cleanParsedUrl)
	}

}

// initialCheck performs basic heuristic checks on the possible URL string:
// - 3 < length < 600
// - no tabs-whitespaces
// - no 'mailto:'
func initialCheck(possibleUrl string) error {
	if possibleUrl == "" || len([]rune(possibleUrl)) < 3 || len(possibleUrl) > 600 {
		return fmt.Errorf("invalid_length")
	} else if strings.Contains(possibleUrl, " ") || strings.Contains(possibleUrl, "\n") || strings.HasPrefix(possibleUrl, "mailto") {
		return fmt.Errorf("possibly_not_url")
	} else {
		return nil
	}

}

// cleanAndParseUrl cleans the URL string and parses it, returning (parsed URL, is local bool, error).
func cleanAndParseUrl(possibleUrl string, domain string, registeredDomain string) (*url.URL, bool, error) {
	cleanedUrl := cleanUrl(possibleUrl)
	return solveUrl(cleanedUrl, domain, registeredDomain)
}

// cleanUrl removes fragments (#), wildcards (*) and adds protocol (HTTPs) in case of a
// protocol-relative url (e.g.: "//example.com/path1")
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

// solveUrl parses a cleaned string URL (coming from cleanUrl()), resolves relative URLs, and determines if it is local.
// Returns (parsed url, is local bool, error).
func solveUrl(cleanedUrl string, domain string, registeredDomain string) (*url.URL, bool, error) {

	// solve relative URLs
	if strings.HasPrefix(cleanedUrl, `/`) {

		cleanedUrl = `https://` + domain + cleanedUrl

		cleanedParsedUrl, err := url.Parse(cleanedUrl)
		if err != nil {
			return nil, true, err
		}

		return cleanedParsedUrl, true, nil
	}

	// add protocol if missing in cases like example.com/path1/path2. Protocol-relative URLs (starts with //) are handled in cleanUrl()
	if !strings.HasPrefix(cleanedUrl, "http") {
		cleanedUrl = `https://` + cleanedUrl
	}

	cleanedParsedUrl, err := url.Parse(cleanedUrl)
	if cleanedParsedUrl == nil || err != nil {
		return nil, false, err // here we return false for a possible internal, but who cares
	}

	// extract its registered domain to determine if it's local or external
	foundUrlRegDomain, err := publicsuffix.EffectiveTLDPlusOne(cleanedParsedUrl.Hostname())
	if err != nil {
		return nil, false, err // here we do it too :)
	}

	// now we can determine locality
	if foundUrlRegDomain == registeredDomain {
		return cleanedParsedUrl, true, nil
	} else {
		return cleanedParsedUrl, false, nil
	}

}
