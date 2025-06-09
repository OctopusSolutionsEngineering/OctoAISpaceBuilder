package validation

import (
	"errors"
	"github.com/OctopusSolutionsEngineering/OctoAISpaceBuilder/internal/domain/terraform"
	"os"
	"path/filepath"
)

// TestFileSystemProviderInstallation is a sanity check to ensure the filesystem provider installation is correct.
func TestFileSystemProviderInstallation(baseDir string) error {
	octopusProvidersDir := filepath.Join(
		baseDir,
		"provider",
		"registry.opentofu.org",
		"octopusdeploy",
		"octopusdeploy",
		terraform.TerraformProviderVersion,
		"linux_amd64",
		"terraform-provider-octopusdeploy_v"+terraform.TerraformProviderVersion)

	if _, err := os.Stat(octopusProvidersDir); os.IsNotExist(err) {
		return errors.New("directory " + octopusProvidersDir + " does not exist")
	}

	return nil
}

func TestOpaPolicyInstallation(baseDir string) error {
	opaPolicyPath := filepath.Join(baseDir, "policy", "terraform.rego")

	if _, err := os.Stat(opaPolicyPath); os.IsNotExist(err) {
		return errors.New("OPA policy file " + opaPolicyPath + " does not exist")
	}

	return nil
}
