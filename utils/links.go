package utils

import (
	"net/url"
)

func ResolveUri(base, uri string) (string, error) {
	baseUrl, err := url.Parse(base)
	if err != nil {
		return "", err
	}

	uriUrl, err := url.Parse(uri)
	if err != nil {
		return "", err
	}

	return baseUrl.ResolveReference(uriUrl).String(), nil
}
