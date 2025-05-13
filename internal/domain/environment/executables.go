package environment

import (
	"os"
	"strings"
)

func DisableMakeBinariesExecutable() bool {
	// Get the disable validation flag from the environment variable
	disableBinariesExecutable := os.Getenv("DISABLE_BINARIES_EXECUTABLE")
	if disableBinariesExecutable == "" {
		return false // Default to false if not set
	}
	return strings.ToLower(disableBinariesExecutable) == "true"
}
