package environment

import "os"

func DisableMakeBinariesExecutable() bool {
	// Get the disable validation flag from the environment variable
	disableBinariesExecutable := os.Getenv("DISABLE_BINARIES_EXECUTABLE")
	if disableBinariesExecutable == "" {
		return false // Default to false if not set
	}
	return disableBinariesExecutable == "true"
}
