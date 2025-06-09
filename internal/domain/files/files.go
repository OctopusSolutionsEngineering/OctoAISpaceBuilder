package files

import (
	"os"
	"path/filepath"
)
import cp "github.com/otiai10/copy"

func CreateTempDir() (string, error) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "tempdir")
	if err != nil {
		return "", err
	}

	// Return the path to the temporary directory
	return tempDir, nil
}

func CopyDir(source string) (string, error) {
	if source == "" {
		return "", nil
	}

	dest, err := os.MkdirTemp("", "octoterra")
	if err != nil {
		return "", err
	}
	err = cp.Copy(source, dest)

	return dest, err
}

func GetAbsoluteOrRelativePath(relativeOrAbsolute string, basePath string) string {
	if relativeOrAbsolute == "" {
		return basePath
	}

	// Check if the path is already absolute
	if filepath.IsAbs(relativeOrAbsolute) {
		return relativeOrAbsolute
	}

	// If not, make it relative to the current working directory
	return filepath.Join(basePath, relativeOrAbsolute)
}
