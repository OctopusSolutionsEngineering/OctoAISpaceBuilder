package environment

import (
	"os"
	"path/filepath"
)

func GetTofuExecutable() (string, error) {
	if os.Getenv("SPACEBUILDER_TOFU_PATH") != "" {
		return os.Getenv("SPACEBUILDER_TOFU_PATH"), nil
	}

	path, err := GetInstallationDirectory()

	if err != nil {
		return "", err
	}

	return filepath.Join(path, "binaries/tofu"), nil
}
