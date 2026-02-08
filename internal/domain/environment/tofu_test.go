package environment

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTofuExecutable(t *testing.T) {
	// Save original environment variable to restore later
	originalTofuPath := os.Getenv("SPACEBUILDER_TOFU_PATH")

	// Clean up after the test
	defer func() {
		os.Setenv("SPACEBUILDER_TOFU_PATH", originalTofuPath)
	}()

	tests := []struct {
		name         string
		tofuPath     string
		expectedPath string
	}{
		{
			name:         "uses custom path when environment variable is set",
			tofuPath:     "/custom/path/to/tofu",
			expectedPath: "/custom/path/to/tofu",
		},
		{
			name:         "uses default path when environment variable is not set",
			tofuPath:     "",
			expectedPath: "binaries/tofu",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variable for this test case
			os.Setenv("SPACEBUILDER_TOFU_PATH", tt.tofuPath)

			// Call the function under test
			path, err := GetTofuExecutable()

			if err != nil {
				t.Fatalf("GetTofuExecutable returned an error: %v", err)
			}

			// Assert
			assert.Equal(t, tt.expectedPath, path)
		})
	}
}
