package environment

import (
	"os"
	"strings"
)

func GetDisableTerraformCliConfig() bool {
	// Check if the environment variable is set to "true"
	if strings.ToLower(os.Getenv("SPACEBUILDER_DISABLE_TERRAFORM_CLI_CONFIG")) == "true" {
		return true
	}

	// Default to false if the environment variable is not set or not "true"
	return false
}
