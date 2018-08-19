package adstxt

// Warning represent failure to parse Ads.txt line according to official ads.txt spec
type Warning struct {
	Index   int      `json:"index"` // Index of the line in the Ads.txt file in which warning was found
	Text    string   `json:"txt"`   // Text of the line in the Ads.txt file in which warning was found
	Message string   `json:"msg"`   // Warning reason
	Level   Sevirity `json:"level"` // Sevirity level of parse warning
}

// Sevirity of parse warning (low for moderate warning, high indicates potential erro)
type Sevirity int

const (
	// ignore first value by assigning to blank identifier
	_ = iota
	// LowSevirity sevirity level for parse warning (low)
	LowSevirity Sevirity = 1 + iota // Warning indicate
	// HighSevirity lsevirity level for parse warning (high, indicates possible error)
	HighSevirity
)
