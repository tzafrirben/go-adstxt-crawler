package adstxt

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

// Calling remote host error\warning
const (
	errHTTPClientError    = "[%s] remote host [%s] Ads.txt URL [%s]"
	errHTTPGeneralError   = "[%s] remote host [%s] Ads.txt URL [%s]"
	errHTTPBadContentType = "[%s] Ads.txt file content type should be ‘text/plain’ and not [%s]"
)

// parsing error\warning: each error includes Ads.txt remote host (domain level) and explanaiton about the error
const (
	errFailToParseRedirect       = "[%s] failed to parse root domain from HTTP redirect response header. Ads.txt URL [%s] redirect [%s] error [%s]"
	errRedirectToInvalidAdsTxt   = "[%s] failed to get Ads.txt file, redirect from [%s] to invalid Ads.txt URL [%s]"
	errRedirectToDifferentDomain = "Only single redirect out of original root domain scope [%s] is allowed. Additional redirect from [%s] to [%s] is forbidden"
)

// HTTP crawler settings
const (
	userAgent      = "+https://github.com/tzafrirben/go-adstxt-crawler"
	requestTimeout = 30
)

// crawler provide methods for downloading Ads.txt files from remote host
type crawler struct {
	client    *http.Client // HTTP client used to make HTTP request for Ads.txt file from remote host
	UserAgent string       // crawler UserAgent string
}

// NewCrawler Create new crawler to fetch Ads.txt file from remote host
func newCrawler() *crawler {
	return &crawler{
		// Create client with required custom parameters.
		// Options: Disable keep-alives, 30sec n/w call timeout, do not follow redirects by default
		client: &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
			Transport: &http.Transport{
				DisableKeepAlives: true,
			},
			Timeout: time.Second * requestTimeout,
		},
		UserAgent: userAgent,
	}
}

// send HTTP request to fetch Ads.txt file from remote host
func (c *crawler) sendRequest(req *Request) (*http.Response, error) {
	httpRequest, err := http.NewRequest("GET", req.URL, nil)
	if err != nil {
		return nil, err
	}

	httpRequest.Header.Add("User-Agent", c.UserAgent)
	httpRequest.Header.Add("Accept", "text/plain")
	httpRequest.Header.Add("Accept-Charset", "utf-8")
	httpRequest.Header.Add("Content-Type", "text/plain; charset=utf-8")

	res, err := c.client.Do(httpRequest)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// handle HTTP redirect response: parse new redirect destination from HTTP response header
func (c *crawler) handleRedirect(req *Request, res *http.Response) (string, error) {
	redirect := res.Header.Get("Location")

	log.Printf("[%s]: redirect from [%s] to [%s]", res.Status, req.URL, redirect)

	// Check if redirect destination has the same root domain as the request initial root doamin.
	d, err := rootDomain(redirect)
	if err != nil {
		return "", fmt.Errorf(errFailToParseRedirect, req.Domain, req.URL, redirect, err.Error())
	}

	// According to IAB ads.txt specification, section 3.1 "ACCESS METHOD":
	// "If the server response indicates an HTTP/HTTPS redirect (301, 302, 307 status codes),
	// the advertising system should follow the redirect and consume the data as authoritative for the source of the redirect,
	// if and only if the redirect is within scope of the original root domain as defined above.
	// Multiple redirects are valid as long as each redirect location remains within the original root domain."
	if d != req.Domain {
		// If redirect to different domain, check that this is the first redirect to different domain
		// According to IAB ads.txt specification, section 3.1 "ACCESS METHOD":
		// "Only a single HTTP redirect to a destination outside the original root domain is allowed to
		// facilitate one-hop delegation of authority to a third party's web server domain."
		prevDomain, _ := rootDomain(req.URL)
		if prevDomain != req.Domain && prevDomain != d {
			return "", fmt.Errorf(errRedirectToDifferentDomain, req.Domain, prevDomain, d)
		}
	}

	// make sure redirects takes us to another Ads.txt file and not just to home page
	if !strings.HasSuffix(redirect, "/ads.txt") {
		return "", fmt.Errorf(errRedirectToInvalidAdsTxt, req.Domain, req.URL, redirect)
	}

	return redirect, nil
}

// Read HTTP response body
func (c *crawler) readBody(req *Request, res *http.Response) ([]byte, error) {
	// The HTTP Content-type should be ‘text/plain’, and all other Content-types should be treated as
	// an error and the content ignored
	contentType := res.Header.Get("Content-Type")
	if strings.Index(contentType, "text/plain") != 0 {
		return nil, fmt.Errorf(errHTTPBadContentType, req.URL, contentType)
	}

	// read response body
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// parse Ads.txt file expiration date from the response Expires header
func (c *crawler) parseExpires(res *http.Response) (time.Time, error) {
	expires := res.Header.Get("Expires")
	if len(expires) == 0 {
		return time.Time{}, fmt.Errorf("Failed to parse expires from response header")
	}

	parsedHeader, err := http.ParseTime(expires)
	if err != nil {
		log.Printf("[%s] Error when parsing HTTP expires header from response [%s]", res.Request.URL, err.Error())
		return time.Time{}, err
	}

	return parsedHeader, nil
}
