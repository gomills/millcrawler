package utils

import (
	"net/url"

	"golang.org/x/net/publicsuffix"
)

// IsLocalUrl determines whether the given parsedUrl belongs to the specified registeredDomain.
func IsLocalUrl(parsedUrl *url.URL, registeredDomain string) (bool, error) {
	urlRegDom, err := publicsuffix.EffectiveTLDPlusOne(parsedUrl.Hostname())
	if err != nil {
		return false, err
	}
	return urlRegDom == registeredDomain, nil
}
