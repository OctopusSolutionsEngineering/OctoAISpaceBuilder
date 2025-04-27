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

func TestGenerateOverrides(t *testing.T) {
	// Call the function under test
	result := GenerateOverrides()

	// Verify the content
	assert.Contains(t, result, "required_providers {", "Should contain required_providers section")
	assert.Contains(t, result, "octopusdeploy = { source = \"OctopusDeployLabs/octopusdeploy\", version = \""+TerraformProviderVersion+"\" }",
		"Should contain correct octopusdeploy provider definition")
	assert.Contains(t, result, "required_version = \">= 1.6.0\"", "Should contain minimum Terraform version")

	// Verify the structure
	lines := strings.Split(strings.TrimSpace(result), "\n")
	assert.GreaterOrEqual(t, len(lines), 5, "Should contain at least 5 lines")

	// Verify the provider version matches the constant
	assert.Contains(t, result, "version = \""+TerraformProviderVersion+"\"",
		"Provider version should match TerraformProviderVersion constant")
}

func TestGenerateStateFile(t *testing.T) {
	// Call the function under test
	result := GenerateStateFile()

	// Verify the content
	assert.Contains(t, result, "terraform {", "Should contain terraform block")
	assert.Contains(t, result, "backend \"local\" {", "Should contain local backend block")
	assert.Contains(t, result, "path = \"./.local-state\"", "Should define local state path correctly")

	// Verify structure
	lines := strings.Split(strings.TrimSpace(result), "\n")
	assert.GreaterOrEqual(t, len(lines), 5, "Should contain at least 5 lines")

	// Verify proper formatting
	assert.True(t, strings.Contains(result, "{") && strings.Contains(result, "}"),
		"Should contain properly formed blocks with braces")
}
