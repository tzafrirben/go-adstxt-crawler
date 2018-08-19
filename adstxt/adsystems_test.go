package adstxt

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// add test values
	adSystems[10000] = newAdSystem(10000, "greenadexchange", "greenadexchange.com")
	adSystemDomains["greenadexchange.com"] = newAdSystemDomain("greenadexchange.com", 10000)

	adSystems[10001] = newAdSystem(10001, "testexchange", "testexchange.net")
	adSystemDomains["testexchange.com"] = newAdSystemDomain("testexchange.com", 10001)

	os.Exit(m.Run())
}

// TestValidateDomain test ads.txt validate domain helper function
func TestValidateDomain(t *testing.T) {
	// collection of domains and expected validation result
	domains := map[string]bool{
		"example.com":               true,
		"http://example.com":        false,
		"www.example.com":           true,
		"a.example.co.uk":           true,
		"ads.example.com":           true,
		"example.com/path":          false,
		"rtb.selectmedia.asia":      true,
		"rtb.selectmedia.asia/path": false,
	}

	for k, v := range domains {
		res := validateDomainName(k)
		if res != v {
			t.Errorf("Expected [%s] domain validation to be [%t]", k, v)
		}
	}
}

// TestRootDomain test exporting root domain from specified URL. Root domain is defined as
// the “public suffix” plus one string in the name
func TestRootDomain(t *testing.T) {
	domains := map[string]string{
		"example.com":                    "example.com",
		"http://example.com":             "example.com",
		"https://example.com":            "example.com",
		"https://example.com/":           "example.com",
		"http://www.example.com":         "example.com",
		"www.example.com/":               "example.com",
		"www.example.com/ads.txt":        "example.com",
		"http://www.test.com/ads.txt":    "test.com",
		"example.com/path/":              "example.com",
		"sub-domain.test.com":            "test.com",
		"http://sub.domain.test.com":     "test.com",
		"http://abc.raisingourkids.com/": "raisingourkids.com",
		"https://testme.tumblr.com/":     "tumblr.com",
		"http://port.com:8080/grid":      "port.com",
	}

	for k, v := range domains {
		res, err := rootDomain(k)
		if err != nil {
			t.Error(err)
		}
		if res != v {
			t.Errorf("Expected [%s] root domain to be [%s] and not [%s]", k, v, res)
		}
	}
}

// TestAdSystemCanonicalName test validation for AdSystem canonical name
func TestAdSystemCanonicalName(t *testing.T) {
	// valid Canonical names
	domains := []string{
		"adtech.com",
		"aolcloud.net",
		"Tremorhub.com",
		"btrll.com",
	}

	for _, d := range domains {
		err := vaidateAdSystemCName(d)
		if err != nil {
			t.Error(err)
		}
	}

	// invalid Canonical names
	domains = []string{
		"example.com",
		"test.net",
		"rtb.selectmedia.asia/path",
		"testexchange.com",
	}

	for _, d := range domains {
		err := vaidateAdSystemCName(d)
		if err == nil {
			t.Errorf("%s is a valid AdSystem canonical name", d)
		}
	}
}
