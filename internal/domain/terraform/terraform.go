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

	if err := WriteVariableFile(path); err != nil {
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

// WriteVariableFile defines empty values for variables often used with the provider. We rely on the fact that
// tofu will ignore variable values for variables that are not set in the plan.
func WriteVariableFile(path string) error {
	serverVar := `octopus_apikey = ""
octopus_server = ""`

	filePath := filepath.Join(path, "terraform.tfvars")

	return os.WriteFile(filePath, []byte(serverVar), 0644)
}
