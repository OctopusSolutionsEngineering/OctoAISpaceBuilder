package files

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateTempDir(t *testing.T) {
	// Call the function
	tempDir, err := CreateTempDir()

	// Assert no error occurred
	require.NoError(t, err, "CreateTempDir should not return an error")

	// Assert we got a non-empty path
	assert.NotEmpty(t, tempDir, "CreateTempDir should return a non-empty path")

	// Check that the directory exists
	info, err := os.Stat(tempDir)
	assert.NoError(t, err, "The created directory should exist")
	assert.True(t, info.IsDir(), "The created path should be a directory")

	// Verify the directory has the expected prefix
	assert.Contains(t, tempDir, "tempdir", "The directory should contain the specified prefix")

	// Clean up after test
	err = os.RemoveAll(tempDir)
	assert.NoError(t, err, "Should be able to remove the temporary directory")
}
