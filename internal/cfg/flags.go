package cfg

import (
	"os"
	"strconv"
)

var (
	DebugMode             = false
	DebugModeRenderColour = uint32(0xff00ff)
	ScaleFactor           = 1.0
)

func init() {
	val := os.Getenv("DEBUG_MODE")
	if debugMode, _ := strconv.ParseBool(val); debugMode {
		DebugMode = debugMode
	}
}
