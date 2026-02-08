package environment

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetOpaExecutable(t *testing.T) {
	// Save original environment variable to restore later
	originalOpaPath := os.Getenv("SPACEBUILDER_OPA_PATH")

	cwd, err := os.Getwd()

	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}

	// Clean up after the test
	defer func() {
		os.Setenv("SPACEBUILDER_OPA_PATH", originalOpaPath)
	}()

	tests := []struct {
		name         string
		opaPath      string
		expectedPath string
	}{
		{
			name:         "uses custom path when environment variable is set",
			opaPath:      "/custom/path/to/opa",
			expectedPath: "/custom/path/to/opa",
		},
		{
			name:         "uses default path when environment variable is not set",
			opaPath:      "",
			expectedPath: filepath.Join(cwd, "functions/binaries/opa_linux_amd64"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variable for this test case
			os.Setenv("SPACEBUILDER_OPA_PATH", tt.opaPath)

			// Call the function under test
			path, err := GetOpaExecutable()

			if err != nil {
				t.Fatalf("GetOpaExecutable returned an error: %v", err)
			}

			// Assert
			assert.Equal(t, tt.expectedPath, path)
		})
	}
}

func TestGetOpaPolicyPath(t *testing.T) {
	// Save original environment variable to restore later
	originalPolicyPath := os.Getenv("SPACEBUILDER_OPA_POLICY_PATH")

	// Clean up after the test
	defer func() {
		os.Setenv("SPACEBUILDER_OPA_POLICY_PATH", originalPolicyPath)
	}()

	tests := []struct {
		name         string
		policyPath   string
		expectedPath string
	}{
		{
			name:         "uses custom policy path when environment variable is set",
			policyPath:   "/custom/policies",
			expectedPath: "/custom/policies",
		},
		{
			name:         "uses default policy path when environment variable is not set",
			policyPath:   "",
			expectedPath: "policy/",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variable for this test case
			os.Setenv("SPACEBUILDER_OPA_POLICY_PATH", tt.policyPath)

			// Call the function under test
			path := GetOpaPolicyPath()

			// Assert
			assert.Equal(t, tt.expectedPath, path)
		})
	}
}
