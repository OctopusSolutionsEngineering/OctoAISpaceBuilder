package environment

import (
	"os"
	"testing"
)

func TestDisableMakeBinariesExecutable(t *testing.T) {
	// Save original environment variable to restore later
	originalDisableBinariesExecutable := os.Getenv("DISABLE_BINARIES_EXECUTABLE")

	// Clean up after the test
	defer func() {
		os.Setenv("DISABLE_BINARIES_EXECUTABLE", originalDisableBinariesExecutable)
	}()

	tests := []struct {
		name         string
		envValue     string
		expectedBool bool
	}{
		{
			name:         "returns true when env var is 'true'",
			envValue:     "true",
			expectedBool: true,
		},
		{
			name:         "returns false when env var is not 'true'",
			envValue:     "false",
			expectedBool: false,
		},
		{
			name:         "returns false when env var is empty",
			envValue:     "",
			expectedBool: false,
		},
		{
			name:         "returns false when env var is set to another value",
			envValue:     "1",
			expectedBool: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up environment for this test case
			os.Setenv("DISABLE_BINARIES_EXECUTABLE", tt.envValue)

			// Call the function under test
			result := DisableMakeBinariesExecutable()

			// Assert
			if tt.expectedBool != result {
				t.Fatalf("DisableMakeBinariesExecutable() expected %v, got %v", tt.expectedBool, result)
			}
		})
	}
}
