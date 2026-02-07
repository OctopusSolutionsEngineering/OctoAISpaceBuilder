package environment

import (
	"os"

	"github.com/OctopusSolutionsEngineering/OctoAISpaceBuilder/internal/domain/files"
)

func GetOpaExecutable() (string, error) {
	if os.Getenv("SPACEBUILDER_OPA_PATH") != "" {
		return os.Getenv("SPACEBUILDER_OPA_PATH"), nil
	}

	installDir, err := GetInstallationDirectory()

	if err != nil {
		return "", err
	}

	return files.GetAbsoluteOrRelativePath(installDir, "binaries/opa_linux_amd64"), nil
}

func GetOpaPolicyPath() string {
	if os.Getenv("SPACEBUILDER_OPA_POLICY_PATH") != "" {
		return os.Getenv("SPACEBUILDER_OPA_POLICY_PATH")
	}

	return "policy/"
}

func GetCombinedOpaPolicyPath() (string, error) {
	installDir, err := GetInstallationDirectory()

	if err != nil {
		return "", err
	}

	policyPath := GetOpaPolicyPath()

	return files.GetAbsoluteOrRelativePath(policyPath, installDir), nil
}
