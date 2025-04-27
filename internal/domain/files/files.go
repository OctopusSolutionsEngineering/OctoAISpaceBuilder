package files

import "os"
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
