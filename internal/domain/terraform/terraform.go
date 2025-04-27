package terraform

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const TerraformProviderVersion = "0.42.0"

func WriteOverrides(path string) error {
	if err := WriteBackendOverride(path); err != nil {
		return err
	}

	if err := WriteProviderOverrides(path); err != nil {
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
		octopusdeploy = { source = "OctopusDeployLabs/octopusdeploy", version = "` + TerraformProviderVersion + `" }
	  }
	  required_version = ">= 1.6.0"
	}`

	filePath := filepath.Join(path, "provider_override.tf")

	return os.WriteFile(filePath, []byte(providerOverrides), 0644)
}

// createTerraformRcFile creates a .terraformrc file in the user's home directory
// The providers directory structure needs to be like:
// provider/registry.terraform.io/octopusdeploylabs/octopusdeploy/0.41.0/linux_amd64/terraform-provider-octopusdeploy_v0.41.0
func CreateTerraformRcFile() error {
	content, err := GenerateTerraformRC()

	if err != nil {
		return fmt.Errorf("failed to generate terraform rc file content: %w", err)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to determine user home directory: %w", err)
	}

	rcFilePath := filepath.Join(homeDir, ".terraformrc")

	// Check if file already exists
	if err := BackupRcFile(rcFilePath); err != nil {
		return err
	}

	// Write the new content
	if err := os.WriteFile(rcFilePath, []byte(content), 0600); err != nil {
		return fmt.Errorf("failed to write .terraformrc file: %w", err)
	}

	return nil
}

func GenerateTerraformRC() (string, error) {
	currentDir, err := os.Getwd()

	if err != nil {
		return "", err
	}

	return `provider_installation {
  filesystem_mirror {
    path    = "` + currentDir + `/provider"
    include = ["*/*/*"]
  }
  direct {}
}`, nil

}

func BackupRcFile(rcFilePath string) error {
	if _, err := os.Stat(rcFilePath); err == nil {
		// Back up existing file
		backupPath := rcFilePath + ".backup." + time.Now().Format("20060102150405")
		if err := os.Rename(rcFilePath, backupPath); err != nil {
			return fmt.Errorf("failed to backup existing .terraformrc file: %w", err)
		}
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("failed to check if .terraformrc file exists: %w", err)
	}

	return nil
}
