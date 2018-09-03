package adstxt

import (
	"fmt"
	"net/url"
	"strings"
)

// Request to fetch Ads.txt file from remote host
type Request struct {
	Domain string `json:"domain"` // Domain holds the root domain of the remote host
	URL    string `json:"url"`    // URL of the Ads.txt file to fetch
}

// NewRequest create new Ads.txt file request from remote host
func NewRequest(rawurl string) (*Request, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}

	// add scheme to Ads.txt URL if it's missing (by default we will add http and not https since it seems more common. If the site is
	// running using HTTPS, we will usually get an HTTP redirect response and will handle it)
	if u.Scheme == "" {
		u.Scheme = "http"
	}

	// add "/ads.txt" to URL path
	if !strings.HasSuffix(u.Path, "/ads.txt") {
		u.Path = fmt.Sprintf("%s/ads.txt", strings.TrimSuffix(u.Path, "/"))
	}

	// Publishers should post the "/ads.txt" file on their root domain and any subdomains as needed.
	// Root domain is defined as the “public suffix” plus one string in the name
	d, err := rootDomain(rawurl)
	if err != nil {
		return nil, err
	}

	adsTxtURL := fmt.Sprintf("%v", u)
	return &Request{URL: adsTxtURL, Domain: d}, nil
}
