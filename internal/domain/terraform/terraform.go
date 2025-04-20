package terraform

import (
	"os"
	"path/filepath"
)

func WriteOverrides(path string) error {
	if err := WriteBackendOverride(path); err != nil {
		return err
	}

	if err := WriteProviderOverrides(path); err != nil {
		return err
	}

	if err := WriteProviderServerVariableOverrides(path); err != nil {
		return err
	}

	if err := WriteProviderApiKeyVariableOverrides(path); err != nil {
		return err
	}

	return nil
}

// WriteBackendOverride forces the use of local state. We never want to send state to the cloud.
func WriteBackendOverride(path string) error {
	backendOverride := `terraform {
	  backend "local" {
		path = "./.local-state"
	  }
	}`

	filePath := filepath.Join(path, "backend_override.tf")

	return os.WriteFile(filePath, []byte(backendOverride), 0644)
}

// WriteProviderOverrides defines the provider block for the Octopus Deploy provider.
// It sets the version to the bundled provider and removes the optional providers that are not needed.
func WriteProviderOverrides(path string) error {
	providerOverrides := `terraform {
	  required_providers {
		octopusdeploy = { source = "OctopusDeployLabs/octopusdeploy", version = "0.42.0" }
	  }
	  required_version = ">= 1.6.0"
	}`

	filePath := filepath.Join(path, "provider_override.tf")

	return os.WriteFile(filePath, []byte(providerOverrides), 0644)
}

// WriteProviderServerVariableOverrides overrides the space variable often used to define the server in the provider block.
// This value must be a blank string for the provider to use the environment variable. It also requires a default value
// to avoid errors.
func WriteProviderServerVariableOverrides(path string) error {
	serverVar := `variable "octopus_server" {
	  default = ""
	}`

	filePath := filepath.Join(path, "provider_server_variable_override.tf")

	return os.WriteFile(filePath, []byte(serverVar), 0644)
}

// WriteProviderApiKeyVariableOverrides overrides the api key variable often used to define the api key in the provider block.
// This value must be a blank string for the provider to use the environment variable. It also requires a default value
// to avoid errors.
func WriteProviderApiKeyVariableOverrides(path string) error {
	apikeyVar := `variable "octopus_apikey" {
	  default = ""
	}`

	filePath := filepath.Join(path, "provider_api_key_variable_override.tf")

	return os.WriteFile(filePath, []byte(apikeyVar), 0644)
}
