package cfg

import (
	"os"
	"strconv"
)

var (
	DebugMode   = false
	ScaleFactor = 1.0

	// TODO this should be set from somewhere.
	FontFolder = "/Library/Fonts"
)

func init() {
	val := os.Getenv("DEBUG_MODE")
	if debugMode, _ := strconv.ParseBool(val); debugMode {
		DebugMode = debugMode
	}
}
