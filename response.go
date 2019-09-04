package adstxt

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Records holds collection of Ads.txt records parsed from an Ads.txt file, in addition to
// errors found during Ads.txt file parsing
type Records struct {
	DataRecords []*DataRecord `json:"dataRecords"`
	Variables   []*Variable   `json:"variables"`
	Warnings    []*Warning    `json:"warnings"`
	Body        []string      `json:"body"` // Original Ads.txt file content
}

// Response to an Ads.txt request: collection of Data\Variable records parsed from Ads.txt file and
// file expiration date
type Response struct {
	*Request
	*Records
	Expires time.Time `json:"expires"` // Ads.txt file expiration date
}

// parseRecords parse Ads.txt file content
func parseRecords(lines []string) *Records {
	r := &Records{
		DataRecords: []*DataRecord{},
		Variables:   []*Variable{},
		Warnings:    []*Warning{},
		Body:        lines,
	}

	// loop over Ads.txt file lines and parse each line into Ads.txt record
	for index, l := range lines {
		r.parseRecord(index+1, l)
	}

	return r
}

// parseRecord parse a single Ads.txt line into Data\Variable record
func (r *Records) parseRecord(index int, txt string) {
	line := removeComment(txt)

	// ignore comments and empty line
	if len(line) == 0 || string(line) == commentDenote {
		return
	}

	// parse line into Data\Variable record
	if strings.Count(line, ",") >= 2 && strings.Count(line, "=") <= 5 {
		dr, w := parseDataRecord(line)
		if w != nil {
			w.Index = index
			w.Text = txt
			r.Warnings = append(r.Warnings, w)
		}
		if dr != nil {
			r.DataRecords = append(r.DataRecords, dr)
		}
	} else if strings.Index(line, "=") != -1 && strings.Count(line, "=") == 1 {
		v, w := parseVariable(txt)
		if w != nil {
			w.Index = index
			w.Text = txt
			r.Warnings = append(r.Warnings, w)
		} else {
			r.Variables = append(r.Variables, v)
		}
	} else {
		w := &Warning{Text: txt, Index: index, Level: HighSeverity, Message: "could not parse this line"}
		r.Warnings = append(r.Warnings, w)
	}
}

// custom "toString" method
func (r *Records) String() string {
	str := []string{}
	if len(r.Warnings) > 0 {
		str = append(str, fmt.Sprintf("Warnings: [%d]", len(r.Warnings)))
		for _, w := range r.Warnings {
			j, _ := json.Marshal(w)
			str = append(str, string(j))
		}
	}

	if len(r.DataRecords) > 0 {
		str = append(str, fmt.Sprintf("Data Records: [%d]", len(r.DataRecords)))
		for _, r := range r.DataRecords {
			j, _ := json.Marshal(r)
			str = append(str, string(j))
		}
	}

	if len(r.Variables) > 0 {
		str = append(str, fmt.Sprintf("Variables: [%d]", len(r.Variables)))
		for _, v := range r.Variables {
			j, _ := json.Marshal(v)
			str = append(str, string(j))
		}
	}
	return strings.Join(str, "\n")
}
