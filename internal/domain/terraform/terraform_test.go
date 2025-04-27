package terraform

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"strings"
	"testing"
)

func TestGenerateTerraformRC(t *testing.T) {
	// Get current working directory to use in verification
	currentDir, err := os.Getwd()
	require.NoError(t, err, "Failed to get current working directory")

	// Call the function under test
	content, err := GenerateTerraformRC()

	// Verify no error occurred
	require.NoError(t, err, "GenerateTerraformRC should not return an error")

	// Verify the content
	assert.Contains(t, content, "provider_installation {", "Should contain provider_installation section")
	assert.Contains(t, content, "filesystem_mirror {", "Should contain filesystem_mirror section")
	assert.Contains(t, content, "path    = \""+currentDir+"/provider\"", "Should contain correct provider path")
	assert.Contains(t, content, "include = [\"*/*/*\"]", "Should include all providers")
	assert.Contains(t, content, "direct {}", "Should contain direct provider section")

	// Verify overall structure
	lines := strings.Split(strings.TrimSpace(content), "\n")
	assert.GreaterOrEqual(t, len(lines), 5, "Should contain at least 5 lines")
}
