package environment

import (
	"github.com/OctopusSolutionsEngineering/OctoAISpaceBuilder/internal/domain/files"
	"os"
)

func GetOpaExecutable() string {
	if os.Getenv("SPACEBUILDER_OPA_PATH") != "" {
		return os.Getenv("SPACEBUILDER_OPA_PATH")
	}

	return "binaries/opa_linux_amd64"
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
