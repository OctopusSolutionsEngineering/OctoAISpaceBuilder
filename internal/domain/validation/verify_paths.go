package validation

import (
	"errors"
	"github.com/OctopusSolutionsEngineering/OctoAISpaceBuilder/internal/domain/environment"
	"github.com/OctopusSolutionsEngineering/OctoAISpaceBuilder/internal/domain/files"
	"github.com/OctopusSolutionsEngineering/OctoAISpaceBuilder/internal/domain/terraform"
	"os"
	"path/filepath"
)

// TestFileSystemProviderInstallation is a sanity check to ensure the filesystem provider installation is correct.
func TestFileSystemProviderInstallation(baseDir string) error {
	localPath := files.GetAbsoluteOrRelativePath(environment.GetTerraformProvidersPath(), baseDir)

	octopusProvidersDir := filepath.Join(
		localPath,
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
	localPath := files.GetAbsoluteOrRelativePath(environment.GetOpaPolicyPath(), baseDir)

	opaPolicyPath := filepath.Join(localPath, "terraform.rego")

	if _, err := os.Stat(opaPolicyPath); os.IsNotExist(err) {
		return errors.New("OPA policy file " + opaPolicyPath + " does not exist")
	}

	return nil
}
