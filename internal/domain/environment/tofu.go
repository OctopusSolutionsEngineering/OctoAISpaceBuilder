package environment

import (
	"os"

	"github.com/OctopusSolutionsEngineering/OctoAISpaceBuilder/internal/domain/files"
)

func GetTofuExecutable() (string, error) {
	if os.Getenv("SPACEBUILDER_TOFU_PATH") != "" {
		return os.Getenv("SPACEBUILDER_TOFU_PATH"), nil
	}

	path, err := GetInstallationDirectory()

	if err != nil {
		return "", err
	}

	return files.GetAbsoluteOrRelativePath(path, "binaries/tofu"), nil
}
