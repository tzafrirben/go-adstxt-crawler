package adstxt

import (
	"encoding/json"
	"testing"
)

// TestParseDataRecordWith3Fields test parsing Ads.txt data record line with 3 fields
func TestParseDataRecordWith3Fields(t *testing.T) {
	line := "greenadexchange.com,XF7342,  	DIRECT"

	expected := DataRecord{
		AdverterDomain:     "greenadexchange.com",
		PublisherAccountID: "XF7342",
		AccountType:        "DIRECT",
	}

	r, w := parseDataRecord(line)
	if w != nil {
		t.Errorf("Expected no parse warning when parsing [%s] [%v]", line, w)
	}
	if r.AdverterDomain != expected.AdverterDomain {
		t.Errorf("Expected Adverter Domain for [%s] to be [%s] but recieved [%s]", line, expected.AdverterDomain, r.AdverterDomain)
	}
	if r.PublisherAccountID != expected.PublisherAccountID {
		t.Errorf("Expected Publisher Account Id for [%s] to be [%s] but recieved [%s]", line, expected.PublisherAccountID, r.PublisherAccountID)
	}
	if r.AccountType != expected.AccountType {
		t.Errorf("Expected Account Type for [%s] to be [%s] but recieved [%s]", line, expected.AccountType, r.AccountType)
	}
}

// TestParseDateRecordWith4Fields test parsing Ads.txt data record line with 4 fields
func TestParseDateRecordWith4Fields(t *testing.T) {
	line := "greenadexchange.com, XF7342, DIRECT, 5jyxf8k54"

	expected := DataRecord{
		AdverterDomain:     "greenadexchange.com",
		PublisherAccountID: "XF7342",
		AccountType:        "DIRECT",
		CertAuthorityID:    "5jyxf8k54",
	}

	r, w := parseDataRecord(line)
	if w != nil {
		t.Errorf("Expected no parse warnings when parsing [%s] [%v]", line, w)
	}
	if r.AdverterDomain != expected.AdverterDomain {
		t.Errorf("Expected Adverter Domain for [%s] to be [%s] but recieved [%s]", line, expected.AdverterDomain, r.AdverterDomain)
	}
	if r.PublisherAccountID != expected.PublisherAccountID {
		t.Errorf("Expected Publisher Account Id for [%s] to be [%s] but recieved [%s]", line, expected.PublisherAccountID, r.PublisherAccountID)
	}
	if r.AccountType != expected.AccountType {
		t.Errorf("Expected Account Type for [%s] to be [%s] but recieved [%s]", line, expected.AccountType, r.AccountType)
	}
	if r.CertAuthorityID != expected.CertAuthorityID {
		t.Errorf("Expected Account Type for [%s] to be [%s] but recieved [%s]", line, expected.CertAuthorityID, r.CertAuthorityID)
	}
}

// TestParseDataRecordWithWrongNumberOfFields test parsing Ads.txt data record line with wrong number of fields
func TestParseDataRecordWithWrongNumberOfFields(t *testing.T) {
	line := "greenadexchange.com, XF7342"

	r, e := parseDataRecord(line)
	if e == nil {
		t.Errorf("Expected error when parsing [%s] [%v]", line, r)
	}

	line = "greenadexchange.com, XF7342, DIRECT, 5jyxf8k54, not-valid"

	r, e = parseDataRecord(line)
	if e == nil {
		t.Errorf("Expected error when parsing [%s] [%v]", line, r)
	}
}

// TestParseDataRecordAccountType test parsing DataRecord account type field
func TestParseDataRecordAccountType(t *testing.T) {
	// invalid accont type
	line := "greenadexchange.com, XF7342, unknow"

	r, w := parseDataRecord(line)
	if w == nil {
		t.Errorf("Expected parse warnings when parsing [%s] [%v]", line, w)
	}

	// case sensative account type
	line = "greenadexchange.com, XF7342, direct"
	r, w = parseDataRecord(line)
	if w != nil {
		t.Errorf("Expected no parse warnings when parsing [%s] [%v]", line, w)
	}

	if r.AccountType != accountTypeReseller && r.AccountType != accountTypeDirect {
		t.Errorf("Expected account type [%s] to be a valid account type", r.AccountType)
	}
}

// TestParseDataRecordWithComment test parsing of DataRecord line with comment
func TestParseDataRecordWithComment(t *testing.T) {
	line := removeComment("advertising.com,17429, DIRECT, #video, US")

	expected := DataRecord{
		AdverterDomain:     "advertising.com",
		PublisherAccountID: "17429",
		AccountType:        "DIRECT",
	}

	r, w := parseDataRecord(line)
	if w != nil {
		t.Errorf("Expected no errors when parsing [%s] [%v]", line, w)
	}

	if r.AdverterDomain != expected.AdverterDomain {
		t.Errorf("Expected Adverter Domain for [%s] to be [%s] but recieved [%s]", line, expected.AdverterDomain, r.AdverterDomain)
	}
	if r.PublisherAccountID != expected.PublisherAccountID {
		t.Errorf("Expected Publisher Account Id for [%s] to be [%s] but recieved [%s]", line, expected.PublisherAccountID, r.PublisherAccountID)
	}
	if r.AccountType != expected.AccountType {
		t.Errorf("Expected Account Type for [%s] to be [%s] but recieved [%s]", line, expected.AccountType, r.AccountType)
	}
}

// TestParseDataRecordWithInvalidAdvertisingSystemName test parsing Ads.txt data record line with in valid domain name of the advertising system
func TestParseDataRecordWithInvalidAdvertisingSystemName(t *testing.T) {
	line := "greenadexchange,XF7342,DIRECT"
	_, w := parseDataRecord(line)
	if w == nil {
		t.Errorf("Expected error when parsing [%s] [%v]", line, w)
	}

	line = "greenadexchange.com, 185, RESELLER"
	r, err := parseDataRecord(line)
	if err != nil {
		t.Errorf("Expected no error when parsing [%s]", line)
	}

	expected := DataRecord{
		AdverterDomain:     "greenadexchange.com",
		PublisherAccountID: "185",
		AccountType:        "RESELLER",
	}

	// compare DataRecords
	if r.AdverterDomain != expected.AdverterDomain || r.PublisherAccountID != expected.PublisherAccountID || r.AccountType != expected.AccountType {
		t.Errorf("Failed to compare parsed DataRecord [%v] to expected DataRecord [%v]", r, expected)
	}
}

// TestParseDataRecordWithInvalidCertName test parsing Ads.txt data record line with in valid dertification authority
func TestParseDataRecordWithInvalidCertName(t *testing.T) {
	line := "greenadexchange.com,185,DIRECT,<invalid>"
	r, w := parseDataRecord(line)

	if w == nil {
		t.Errorf("Expected error when parsing [%s] [%v]", line, w)
	}

	expected := DataRecord{
		AdverterDomain:     "greenadexchange.com",
		PublisherAccountID: "185",
		AccountType:        "DIRECT",
		CertAuthorityID:    "<invalid>",
	}

	// compare DataRecords
	if r.AdverterDomain != expected.AdverterDomain || r.PublisherAccountID != expected.PublisherAccountID || r.AccountType != expected.AccountType {
		t.Errorf("Failed to compare parsed DataRecord [%v] to expected DataRecord [%v]", r, expected)
	}
}

// TestDataRecordJsonEncode test encoding DataRecord to json
func TestDataRecordJsonEncode(t *testing.T) {
	line := "greenadexchange.com, XF7342, DIRECT, 5jyxf8k54"

	r, w := parseDataRecord(line)
	if w != nil {
		t.Errorf("Expected no errors when parsing [%s] [%v]", line, w)
	}

	j, _ := json.Marshal(r)
	if string(j) != "{\"adverterdomain\":\"greenadexchange.com\",\"publisheraccountid\":\"XF7342\",\"accountype\":\"DIRECT\",\"certauthorityid\":\"5jyxf8k54\"}" {
		t.Errorf("Json encoded DataRecord is different than expected [%s]", string(j))
	}

	// test Json encode without optional value
	line = "greenadexchange.com, XF7342, DIRECT"
	r, w = parseDataRecord(line)
	if w != nil {
		t.Errorf("Expected no errors when parsing [%s] [%v]", line, w)
	}

	j, _ = json.Marshal(r)
	if string(j) != "{\"adverterdomain\":\"greenadexchange.com\",\"publisheraccountid\":\"XF7342\",\"accountype\":\"DIRECT\"}" {
		t.Errorf("Json encoded DataRecord is different than expected [%s]", string(j))
	}
}

// TestDataRecordJsonDecode test decode DataRecord from json
func TestDataRecordJsonDecode(t *testing.T) {
	j := []byte("{\"adverterdomain\":\"greenadexchange.com\",\"publisheraccountid\":\"XF7342\",\"accountype\":\"DIRECT\",\"certauthorityid\":\"5jyxf8k54\"}")

	var r DataRecord
	err := json.Unmarshal(j, &r)
	if err != nil {
		t.Error(err)
	}

	expected := DataRecord{
		AdverterDomain:     "greenadexchange.com",
		PublisherAccountID: "XF7342",
		AccountType:        "DIRECT",
		CertAuthorityID:    "5jyxf8k54",
	}

	if r.AdverterDomain != expected.AdverterDomain {
		t.Errorf("Expected Adverter Domain to be [%s] but recieved [%s]", expected.AdverterDomain, r.AdverterDomain)
	}
	if r.PublisherAccountID != expected.PublisherAccountID {
		t.Errorf("Expected Publisher Account Id to be [%s] but recieved [%s]", expected.PublisherAccountID, r.PublisherAccountID)
	}
	if r.AccountType != expected.AccountType {
		t.Errorf("Expected Account Type to be [%s] but recieved [%s]", expected.AccountType, r.AccountType)
	}
	if r.CertAuthorityID != expected.CertAuthorityID {
		t.Errorf("Expected Account Type to be [%s] but recieved [%s]", expected.CertAuthorityID, r.CertAuthorityID)
	}
}

// TestParseSubDomainVarialbe test parsing subdomain Variable type
func TestParseSubDomainVarialbe(t *testing.T) {
	subdomain := "subdomain=dev.example.com"

	v, w := parseVarialbe(subdomain)
	if w != nil {
		t.Errorf("Expected no errors when parsing [%s] [%v]", subdomain, w)
	}
	if v.Type != varTypeSubdomain {
		t.Errorf("Expected variable type for [%s] to be [%s] but recieved [%s]", subdomain, varTypeSubdomain, v.Type)
	}
	if v.Value != "dev.example.com" {
		t.Errorf("Expected variable value for [%s] to be [dev.example.com] but recieved [%s]", subdomain, v.Value)
	}
}

// TestParseContactVarialbe test parsing contact Variable type
func TestParseContactVarialbe(t *testing.T) {
	contact := "contact=tzafrir@example.com"

	v, w := parseVarialbe(contact)
	if w != nil {
		t.Errorf("Expected no errors when parsing [%s] [%v]", contact, w)
	}
	if v.Type != varTypeContact {
		t.Errorf("Expected variable type for [%s] to be [%s] but recieved [%s]", contact, varTypeContact, v.Type)
	}
	if v.Value != "tzafrir@example.com" {
		t.Errorf("Expected variable value for [%s] to be [tzafrir@example.com] but recieved [%s]", contact, v.Value)
	}
}

// TestParseNotSupportedVarialble test parsing not supported Varialble type
func TestParseNotSupportedVarialble(t *testing.T) {
	notSupported := "notSupported=dev.example.com"

	v, w := parseVarialbe(notSupported)
	if w == nil {
		t.Errorf("Expected parsing error when parsing [%s]", notSupported)
	}
	if v != nil {
		t.Errorf("Expected parsing error when parsing [%s] and empty variable [%v]", notSupported, v)
	}
	if w.Message != "[notSupported] is not a valid Variable type" {
		t.Errorf("Expected error type for [%s] to be [%s] but recieved [%v]", notSupported, "[notSupported] is not a valid Variable type", w)
	}
}

// TestRemoveComment test creating new line with Ads.txt comment
func TestRemoveComment(t *testing.T) {
	s := "advertising.com,17429, DIRECT, #video, US"

	l := removeComment(s)

	if l == s {
		t.Errorf("Expected comment to be removed from line [%s]", s)
	}

	if l != "advertising.com,17429, DIRECT," {
		t.Errorf("Expected parsed line text to be [%s] and not [%s]", "advertising.com,17429, DIRECT, ", l)
	}

}
