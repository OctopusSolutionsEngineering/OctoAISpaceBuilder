package environment

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPort(t *testing.T) {
	// Save original environment variables to restore later
	originalFunctionsPort := os.Getenv("FUNCTIONS_CUSTOMHANDLER_PORT")
	originalSpacebuilderPort := os.Getenv("SPACEBUILDER_FUNCTIONS_CUSTOMHANDLER_PORT")

	// Clean up after the test
	defer func() {
		os.Setenv("FUNCTIONS_CUSTOMHANDLER_PORT", originalFunctionsPort)
		os.Setenv("SPACEBUILDER_FUNCTIONS_CUSTOMHANDLER_PORT", originalSpacebuilderPort)
	}()

	tests := []struct {
		name             string
		functionsPort    string
		spacebuilderPort string
		expectedPort     string
	}{
		{
			name:             "uses FUNCTIONS_CUSTOMHANDLER_PORT when set",
			functionsPort:    "9000",
			spacebuilderPort: "8000",
			expectedPort:     "9000",
		},
		{
			name:             "falls back to SPACEBUILDER_FUNCTIONS_CUSTOMHANDLER_PORT when FUNCTIONS_CUSTOMHANDLER_PORT not set",
			functionsPort:    "",
			spacebuilderPort: "8000",
			expectedPort:     "8000",
		},
		{
			name:             "uses default port when neither env var is set",
			functionsPort:    "",
			spacebuilderPort: "",
			expectedPort:     "8080",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up environment for this test case
			os.Setenv("FUNCTIONS_CUSTOMHANDLER_PORT", tt.functionsPort)
			os.Setenv("SPACEBUILDER_FUNCTIONS_CUSTOMHANDLER_PORT", tt.spacebuilderPort)

			// Call the function under test
			port := GetPort()

			// Assert
			assert.Equal(t, tt.expectedPort, port)
		})
	}
}
