package adstxt

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestSendAndParseRquest test send HTTP request to remote host to Get Ads.txt file
func TestSendAndParseRquest(t *testing.T) {
	// expected response
	const expected = "greenadexchange.com,XF7342,DIRECT"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, expected)
		w.Header().Set("Content-Type", "text/plain")
	}))
	defer ts.Close()

	// request mock
	req, _ := NewRequest(ts.URL)

	// test send request
	c := newCrawler()
	res, err := c.sendRequest(req)
	if err != nil {
		t.Error(err)
	}

	defer res.Body.Close()

	body, err := c.parseBody(req, res)
	if err != nil {
		t.Error(err)
	}

	if string(body) != expected {
		t.Errorf("Expected response body [%s] to be \"%s\"", string(body), expected)
	}
}

// TestHandleRedirect test crawler handle HTTP redirect response: extract new redirect destination from HTTP resposne
func TestHandleRedirect(t *testing.T) {
	const redirect = "http://gotest.com/ads.txt"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Location", redirect)
		w.WriteHeader(http.StatusMovedPermanently)
	}))
	defer ts.Close()

	// request mock
	req, _ := NewRequest(ts.URL)

	// test send request
	c := newCrawler()
	res, err := c.sendRequest(req)
	if err != nil {
		t.Error(err)
	}

	defer res.Body.Close()

	// parse redirect location
	r, err := c.handleRedirect(req, res)
	if err != nil {
		t.Error(err)
	}

	if r != redirect {
		t.Errorf("Expected redirect destination to be [%s] and not [%s]", redirect, r)
	}
}

// TestParseExpires test parse Ads.txt file expires from HTTP response Header
func TestParseExpires(t *testing.T) {
	// expected response
	var cache = time.Now().AddDate(60, 0, 0)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Expires", cache.Format(http.TimeFormat))
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	// request mock
	req, _ := NewRequest(ts.URL)

	// test send request
	c := newCrawler()
	res, err := c.sendRequest(req)
	if err != nil {
		t.Error(err)
	}

	defer res.Body.Close()

	expires, err := c.parseExpires(res)
	if err != nil {
		t.Error(err)
	}

	if expires.Format(http.TimeFormat) != cache.Format(http.TimeFormat) {
		t.Errorf("After Expected expires [%s] to be [%s]", expires.Format(http.TimeFormat), cache.Format(http.TimeFormat))
	}

}
