package environment

import (
	"encoding/json"
	"go.uber.org/zap"
	"os"
	"path/filepath"
	"strings"
)

func GetPort() string {
	// Get the port from the environment variable
	port := os.Getenv("FUNCTIONS_CUSTOMHANDLER_PORT")
	if port == "" {
		port = os.Getenv("SPACEBUILDER_FUNCTIONS_CUSTOMHANDLER_PORT")
		if port == "" {
			port = "8080" // Default port
		}
	}
	return port
}

func IsInAzureFunctions() bool {
	return os.Getenv("FUNCTIONS_WORKER_RUNTIME") == "custom"
}

func GetInstallationDirectory() (string, error) {
	if os.Getenv("SPACEBUILDER_INSTALLATION_DIRECTORY") != "" {
		// If the environment variable is set, use it as the installation directory
		return os.Getenv("SPACEBUILDER_INSTALLATION_DIRECTORY"), nil
	}

	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	if IsInAzureFunctions() {
		// In Azure Functions, everything is in the current working directory
		return cwd, nil
	}

	// When running locally, the policies and other files are in the "functions" directory.
	// This assumes that most people will run this service locally with the working directory set to the root of the project.
	// You can also define absolute paths for the "policy" and "provider" directories in the environment variables
	// SPACEBUILDER_TERRAFORM_PROVIDERS and SPACEBUILDER_OPA_POLICY_PATH.
	// Note that tests will often use the current working directory of their test file. This means tests need to set
	// absolute paths for the policy and provider directories.
	return filepath.Join(cwd, "functions"), nil
}

func DisableValidation() bool {
	// Get the disable validation flag from the environment variable
	disableValidation := os.Getenv("DISABLE_VALIDATION")
	if disableValidation == "" {
		return false // Default to false if not set
	}
	return strings.ToLower(disableValidation) == "true"
}

// GetEnhancedLoggingInstances returns a list of instances for enhanced logging.
// This allows developers to debug their issues without logging prompt responses for anyone else.
func GetEnhancedLoggingInstances() []string {
	instances := []string{}
	instancesJson := os.Getenv("ENHANCED_LOGGING_INSTANCES")
	if instancesJson == "" {
		return []string{} // Default to empty slice if not set
	}

	err := json.Unmarshal([]byte(instancesJson), &instances)
	if err != nil {
		zap.L().Error("Error parsing JSON:", zap.Error(err))
		return []string{}
	}

	return instances
}
