package adstxt

import (
	"testing"
)

func TestNewRequest(t *testing.T) {
	domains := map[string]Request{
		"example.com":                    Request{URL: "http://example.com/ads.txt", Domain: "example.com"},
		"http://example.com":             Request{URL: "http://example.com/ads.txt", Domain: "example.com"},
		"https://example.com":            Request{URL: "https://example.com/ads.txt", Domain: "example.com"},
		"https://example.com/":           Request{URL: "https://example.com/ads.txt", Domain: "example.com"},
		"http://www.example.com":         Request{URL: "http://www.example.com/ads.txt", Domain: "example.com"},
		"www.example.com/":               Request{URL: "http://www.example.com/ads.txt", Domain: "example.com"},
		"www.example.com/ads.txt":        Request{URL: "http://www.example.com/ads.txt", Domain: "example.com"},
		"http://www.test.com/ads.txt":    Request{URL: "http://www.test.com/ads.txt", Domain: "test.com"},
		"example.com/path/":              Request{URL: "http://example.com/path/ads.txt", Domain: "example.com"},
		"sub-domain.test.com":            Request{URL: "http://sub-domain.test.com/ads.txt", Domain: "test.com"},
		"http://sub.domain.test.com":     Request{URL: "http://sub.domain.test.com/ads.txt", Domain: "test.com"},
		"http://abc.raisingourkids.com/": Request{URL: "http://abc.raisingourkids.com/ads.txt", Domain: "raisingourkids.com"}}

	for k, v := range domains {
		r, _ := NewRequest(k)
		if r.URL != v.URL {
			t.Errorf("Expected Ads.txt for [%s] to be [%s] but recieved [%s]", k, v.URL, r.URL)
		}
		if r.Domain != v.Domain {
			t.Errorf("Expected Domain for [%s] to be [%s] but recieved [%s]", k, v.Domain, r.Domain)
		}
	}
}
