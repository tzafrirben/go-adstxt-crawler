package adstxt

// Warning represent failure to parse Ads.txt line according to official ads.txt spec
type Warning struct {
	Index   int      `json:"index"` // Index of the line in the Ads.txt file in which warning was found
	Text    string   `json:"txt"`   // Text of the line in the Ads.txt file in which warning was found
	Message string   `json:"msg"`   // Warning reason
	Level   Severity `json:"level"` // Severity level of parse warning
}

// Severity of parse warning (low for moderate warning, high indicates potential error)
type Severity int

const (
	// ignore first value by assigning to blank identifier
	_ = iota
	// LowSeverity severity level for parse warning (low)
	LowSeverity Severity = 1 + iota // Warning indicate
	// HighSeverity severity level for parse warning (high, indicates possible error)
	HighSeverity
)
