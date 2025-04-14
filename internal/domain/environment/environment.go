package environment

import "os"

func GetPort() string {
	// Get the port from the environment variable
	port := os.Getenv("FUNCTIONS_CUSTOMHANDLER_PORT")
	if port == "" {
		port = "8080" // Default port
	}
	return port
}

func DisableValidation() bool {
	// Get the disable validation flag from the environment variable
	disableValidation := os.Getenv("DISABLE_VALIDATION")
	if disableValidation == "" {
		return false // Default to false if not set
	}
	return disableValidation == "true"
}
