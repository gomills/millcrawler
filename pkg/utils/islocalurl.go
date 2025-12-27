package utils

import (
	"net/url"

	"golang.org/x/net/publicsuffix"
)

// IsLocalUrl receives a parsed url and determines if it's local to the given registered domain
func IsLocalUrl(parsedUrl *url.URL, registeredDomain string) (bool, error) {

	// get rid of the subdomain, keeping only the effective TLD+1
	urlRegDom, err := publicsuffix.EffectiveTLDPlusOne(parsedUrl.Hostname())
	if err != nil {
		return false, err
	}

	return urlRegDom == registeredDomain, nil
}
