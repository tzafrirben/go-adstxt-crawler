package adstxt

import (
	"strings"
	"testing"
)

// TestParseBody test paring []byte array into []Line array
func TestParseBody(t *testing.T) {
	body := []string{
		"this is the first line",
		"this is the second line",
		"this is the third line",
	}

	b := []byte(strings.Join(body, "\n"))
	res, err := ParseBody(b)

	if err != nil {
		t.Error(err)
	}

	if len(res.Body) != len(body) {
		t.Errorf("Expected number of lines to be %d and not %d", len(body), len(res.Body))
	}

	for index, line := range res.Body {
		if line != body[index] {
			t.Errorf("Expected line #%d to be [%s] and not [%s]", index, body[index], line)
		}
	}

	b = []byte(strings.Join(body, "\r\n"))
	res, err = ParseBody(b)

	if err != nil {
		t.Error(err)
	}

	if len(res.Body) != len(body) {
		t.Errorf("Expected number of lines to be %d and not %d", len(body), len(res.Body))
	}

	for index, line := range res.Body {
		if line != body[index] {
			t.Errorf("Expected line #%d to be [%s] and not [%s]", index, body[index], line)
		}
	}
}

// TestParseRecordsSuccess test parsing to valid Ads.txt file
func TestParseRecordsSuccess(t *testing.T) {
	b := []byte("greenadexchange.com,XF7342,  	DIRECT\ngreenadexchange.com, XF7342, DIRECT, 5jyxf8k54\n#greenadexchange.com,XF7342,DIRECT\nsubdomain=dev.example.com")
	res, err := ParseBody(b)

	if err != nil {
		t.Errorf("Expected no errors [%s]", err.Error())
	}

	if len(res.DataRecords) != 2 {
		t.Errorf("Failed to parse Ads.txt records, expected number of records to be 2 and not [%d]", len(res.DataRecords))
	}

	if len(res.Variables) != 1 {
		t.Errorf("Failed to parse Ads.txt variables, expected number of variables to be 1 and not [%d]", len(res.Variables))
	}

	if len(res.Warnings) > 0 {
		t.Errorf("Expected no warning when parsing lines, but recieved [%d] warnings", len(res.Warnings))
	}
}

// TestParseRecordsFailure test parsing to invalid Ads.txt file
func TestParseRecordsFailure(t *testing.T) {
	b1 := []byte("greenadexchange.com,XF7342,\ngreenadexchange.com, XF7342, DIRECT, 5jyxf8k54\n#greenadexchange.com,XF7342,DIRECT\nsubdomain=dev.example.com")
	res, err := ParseBody(b1)

	if err != nil {
		t.Error(err)
	}

	if len(res.Warnings) == 0 {
		t.Error("Expected warnings when parsing Ads.txt with invalid Ads.txt lines")
	}

	if res.Warnings[0].Message != "Missing type of account/relationship (required)" {
		t.Error("Expected warning message to indicate account type is missing")
	}

	b2 := []byte("###this is a comment\ngreenadexchange.com, XF7342, DIRECT, 5jyxf8k54\n#greenadexchange.com,XF7342,DIRECT\nsubdomains=dev.example.com")
	res, err = ParseBody(b2)

	if len(res.Warnings) == 0 {
		t.Error("Expected warnings when parsing Ads.txt with invalid Variable type line")
	}

	if res.Warnings[0].Message != "[subdomains] is not a valid Variable type" {
		t.Error("Expected warning message to indicate account type is missing")
	}

	if len(res.DataRecords) == 0 {
		t.Error("Expected parsing to valid DataRecords line to success")
	}

}
