package adstxt

import (
	"fmt"
	"regexp"
	"strings"
)

// Ads.txt comment
const (
	// Comment is denoted by the character "#".
	commentDenote = "#"
)

// Ads.txt supported account types
const (
	// Direct indicates that the Publisher (content owner) directly controls the account
	accountTypeDirect = "DIRECT"
	// Reseller indicates that the Publisher has authorized another entity to control the account
	accountTypeReseller = "RESELLER"
)

// Ads.txt supported Variables types
const (
	// Subdomain within the root domain on which Ads.txt can be found
	varTypeSubdomain = "subdomain"
	// Contact information for the owner of the Ads.txt file
	varTypeContact = "contact"
)

// DataRecord hold single Ads.txt data record
type DataRecord struct {
	AdverterDomain     string `json:"adverterdomain"`            // AdverterDomain Domain name of the advertising system (required)
	PublisherAccountID string `json:"publisheraccountid"`        // PublisherAccountID the identifier associated with the seller (required)
	AccountType        string `json:"accountype"`                // AccountType enumeration of the type of account: DIRECT or RESELLER (required)
	CertAuthorityID    string `json:"certauthorityid,omitempty"` // CertAuthorityID An ID that uniquely identifies the advertising system within a certification authority (optional)
}

// Variable hold single of Ads.txt variable record
type Variable struct {
	Type  string `json:"type"`  // Type of variable record. Supported types are subdomain and contact
	Value string `json:"value"` // Value of variable record
}

// parseDataRecord return new DataRecord parsed from single Ads.txt line
func parseDataRecord(line string) (*DataRecord, *Warning) {
	// Data record declaraion: <FIELD #1>, <FIELD #2>, <FIELD #3>, <FIELD #4> (optional)
	fields := strings.Split(line, ",")

	fieldsLen := len(fields)
	if fieldsLen < 3 || fieldsLen > 4 {
		return nil, &Warning{Level: HighSeverity, Message: fmt.Sprintf("Data record must be declared as <FIELD #1>, <FIELD #2>, <FIELD #3>, <FIELD #4> (optional) pattern")}
	}

	// make sure required fields are not empty
	adverterDomain := strings.TrimSpace(fields[0])
	if len(adverterDomain) == 0 {
		return nil, &Warning{Level: HighSeverity, Message: fmt.Sprintf("Missing domain name of the advertising system (required)")}
	}

	if !validateDomainName(adverterDomain) {
		return nil, &Warning{Level: HighSeverity, Message: fmt.Sprintf("%s is not a valid Ad system domain", adverterDomain)}
	}

	// check that advertiser domain is a valid DNS name
	err := vaidateAdSystemCName(adverterDomain)
	if err != nil {
		return nil, &Warning{Level: LowSeverity, Message: err.Error()}
	}

	publisherAccountID := strings.TrimSpace(fields[1])
	if len(publisherAccountID) == 0 {
		return nil, &Warning{Level: HighSeverity, Message: fmt.Sprintf("Missing publisher's Account ID (required)")}
	}

	accountType := strings.TrimSpace(fields[2])
	if len(accountType) == 0 {
		return nil, &Warning{Level: HighSeverity, Message: fmt.Sprintf("Missing type of account/relationship (required)")}
	}

	// make sure account type is supported (case insensitive)
	if strings.ToUpper(accountType) != accountTypeReseller && strings.ToUpper(accountType) != accountTypeDirect {
		return nil, &Warning{Level: HighSeverity, Message: fmt.Sprintf("[%s] is not a valid account type. Account type must be [%s] or [%s]",
			accountType, accountTypeDirect, accountTypeReseller)}
	}

	r := DataRecord{
		AdverterDomain:     adverterDomain,
		PublisherAccountID: publisherAccountID,
		AccountType:        strings.ToUpper(accountType),
	}

	// optional value
	if fieldsLen > 3 {
		certAuthorityID := strings.TrimSpace(fields[3])
		r.CertAuthorityID = certAuthorityID

		// check if cert authority id is alphanumeric (if not, it might indicate an error also it is not part of Ads.txt specification)
		re := regexp.MustCompile("^[a-zA-Z0-9]*$")
		if !re.MatchString(r.CertAuthorityID) {
			return &r, &Warning{
				Level:   LowSeverity,
				Message: fmt.Sprintf("Certification Authority ID %s may not be correct as it is not alphanumeric", r.CertAuthorityID),
			}
		}
	}

	return &r, nil
}

// parseVariable return new Variable record parsed from Ads.txt line
func parseVariable(line string) (*Variable, *Warning) {
	// Variable declaraion: lines in the a pattern of <VARIABLE>=<VALUE>
	fields := strings.Split(line, "=")

	// check that record type is supported, and return new variable of that type
	t := fields[0]
	switch strings.ToLower(t) {
	case varTypeSubdomain:
		return &Variable{
			Type:  varTypeSubdomain,
			Value: fields[1],
		}, nil
	case varTypeContact:
		return &Variable{
			Type:  varTypeContact,
			Value: fields[1],
		}, nil
	default:
		return nil, &Warning{Level: HighSeverity, Message: fmt.Sprintf("[%s] is not a valid Variable type", t)}
	}
}

// removeComment removes any comment from Ads.txt line before parsing
func removeComment(line string) string {
	index := strings.Index(line, commentDenote)
	if index != -1 {
		line = line[0:index]
	}
	return strings.TrimSpace(line)
}
