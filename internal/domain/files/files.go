package files

import "os"

func CreateTempDir() (string, error) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "tempdir")
	if err != nil {
		return "", err
	}

	// Return the path to the temporary directory
	return tempDir, nil
}
