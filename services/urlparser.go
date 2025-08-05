package services

import (
	"errors"
	"fmt"
	"jinovatka/assert"
	"net/netip"
	"net/url"
	"slices"
	"strconv"
	"strings"
)

type UrlParserService struct{}

// This variables can be used to test what the returned error represents and customize user facing message.
// This might need a refactor (as it looks ugly).
// TODO: Add comments explaining each error.
var (
	ErrEmptyUri        = errors.New("the URL is empty")
	ErrForbiddenScheme = errors.New("the URL has forbidden scheme")
	ErrEmptyScheme     = errors.New("the URL has empty scheme")
	ErrLoopback        = errors.New("the URL must not be a loopback adress")
	ErrPrivateIP       = errors.New("the URL host must not be a private IP adress")
	ErrWellKnownPort   = errors.New("the URL port must not be in well-known range")
)

// Parse provided uri. Do some paranoid checks before storing and crawling the URL.
//
// User credentials are removed. Scheme is checked to be http or https
// or filled in if empty and strict is false.
// Host is checked to not contain loopback, private IP or well-known port.
//
// This function does not check the uri lenght.
// It also does not check if the resource itself is malicious or NSFW.
func (service *UrlParserService) ParseAndCleanURL(uri string, strict bool) (*url.URL, error) {
	const (
		http  = "http"
		https = "https"
	)

	// Trim whitespace from malformed user input and remaining carige returns from parsing seed lists.
	uri = strings.TrimSpace(uri)

	// Do we even have something to parse?
	if uri == "" {
		return nil, ErrEmptyUri
	}
	// We could also check that the URI isn't too long, but we leave that decision to the caller.

	// Now parse url using net/url.
	parsedUri, err := url.Parse(uri)
	if err != nil {
		return nil, fmt.Errorf("URL failed to parse: %w", err)
	}

	// First things first, delete any userinfo data.
	// Sending userifo data is probably a mistake or some malicious atempt at something.
	parsedUri.User = nil

	// It is possible that user forgot to put a scheme in the adress.
	// If strict is disabled, then we assume https and we fill it in.
	if !strict && parsedUri.Scheme == "" {
		// We could just do: parsedUri.Scheme = https
		// But the host part of URL was likely parsed as path, so we reparse the URL instead.
		parsedUri, err = url.Parse(https + "://" + uri)
		if err != nil {
			return nil, fmt.Errorf("failed to reparse URL after correcting scheme: %w", err)
		}
	}

	// Only allow http and https schemes.
	allowedSchemes := []string{https, http}
	if !slices.Contains(allowedSchemes, parsedUri.Scheme) {
		if parsedUri.Scheme == "" {
			return nil, ErrEmptyScheme
		}
		return nil, ErrForbiddenScheme
	}

	// IMPORTANT: Never allow loopback and private addresses!
	host := parsedUri.Hostname()
	if host == "localhost" {
		return nil, ErrLoopback
	}
	ip, err := netip.ParseAddr(host)
	// If error is nil, then host is IP adress.
	if err == nil && ip.IsValid() {
		if ip.IsLoopback() {
			return nil, ErrLoopback
		}
		if ip.IsPrivate() {
			return nil, ErrPrivateIP
		}
	}

	// Don't allow ports in range (0 - 1023). There may be some valid reasons
	// for using some of these, but likely it's malicious.
	// https://www.rfc-editor.org/rfc/rfc3986#section-7.2
	portString := parsedUri.Port()
	if portString != "" {
		port, err := strconv.Atoi(portString)
		assert.Must(err == nil, "ParseAndCleanURL: impossible state reached, Atoi(port) returned error")
		if port > 0 && port < 1023 {
			return nil, ErrWellKnownPort
		}
	}

	// If we got here then the URL should be OK.
	return parsedUri, nil
}
