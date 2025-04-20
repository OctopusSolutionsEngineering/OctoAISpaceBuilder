package terraform

import (
	"os"
	"path/filepath"
)

func WriteBackendOverride(path string) error {
	backendOverride := `terraform {
	  backend "local" {
		path = "./.local-state"
	  }
	}`

	filePath := filepath.Join(path, "backend_override.tf")

	return os.WriteFile(filePath, []byte(backendOverride), 0644)
}

// WriteProviderServerVariableOverrides overrides the space variable often used to define the server in the provider block.
// This value must be a blank string for the provider to use the environment variable.
func WriteProviderServerVariableOverrides(path string) error {
	server_var := `variable "octopus_server" {
	  default = ""
	}`

	filePath := filepath.Join(path, "provider_server_variable_override.tf")

	return os.WriteFile(filePath, []byte(server_var), 0644)
}

// WriteProviderApiKeyVariableOverrides overrides the api key variable often used to define the api key in the provider block.
// This value must be a blank string for the provider to use the environment variable.
func WriteProviderApiKeyVariableOverrides(path string) error {
	apikey_var := `variable "octopus_apikey" {
	  default = ""
	}`

	filePath := filepath.Join(path, "provider_api_key_variable_override.tf")

	return os.WriteFile(filePath, []byte(apikey_var), 0644)
}
