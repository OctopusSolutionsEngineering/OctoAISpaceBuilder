package handler

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/OctopusSolutionsEngineering/OctoAISpaceBuilder/internal/domain/execute"
	"github.com/OctopusSolutionsEngineering/OctoAISpaceBuilder/internal/domain/files"
	"github.com/OctopusSolutionsEngineering/OctoAISpaceBuilder/internal/domain/jwt"
	"github.com/OctopusSolutionsEngineering/OctoAISpaceBuilder/internal/domain/model"
	"github.com/OctopusSolutionsEngineering/OctoAISpaceBuilder/internal/domain/sha"
	"github.com/OctopusSolutionsEngineering/OctoAISpaceBuilder/internal/domain/terraform"
	"github.com/OctopusSolutionsEngineering/OctoAISpaceBuilder/internal/infrastructure"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"os"
	"path/filepath"
	"time"
)

func CreateTerraformPlan(token string, terraformInput model.TerraformPlan) (*model.TerraformPlan, error) {

	tempDir, err := files.CreateTempDir()

	if err != nil {
		return nil, err
	}

	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			zap.L().Error("Failed to remove temporary directory", zap.Error(err))
		}
	}()

	if err := os.WriteFile(filepath.Join(tempDir, "terraformInput.tf"), []byte(terraformInput.Configuration), 0644); err != nil {
		return nil, err
	}

	if err := terraform.WriteBackendOverride(tempDir); err != nil {
		return nil, err
	}

	aud, err := jwt.GetJwtAud(token)

	if err != nil {
		return nil, err
	}

	if err := createTerraformRcFile(); err != nil {
		return nil, err
	}

	lockFile, err := initTofu(tempDir)

	if err != nil {
		return nil, err
	}

	planFile, planBinary, err := generatePlan(tempDir, token, aud, terraformInput.SpaceId)

	if err != nil {
		return nil, err
	}

	response := model.TerraformPlan{
		ID:      uuid.New().String(),
		Created: time.Now(),
		Server:  sha.GetSha256Hash(aud),
		SpaceId: terraformInput.SpaceId,
	}

	if err := infrastructure.CreateFeedbackAzureStorageTable(response.ID, planBinary, response.SpaceId, response.Server, lockFile, terraformInput.Configuration); err != nil {
		return nil, err
	}

	planJson, err := generatePlanJson(planFile)

	if err != nil {
		return nil, err
	}

	if err := checkPlan(planJson); err != nil {
		return nil, err
	}

	planText, err := generatePlanText(planFile)

	if err != nil {
		return nil, err
	}

	response.PlanText = &planText

	return &response, nil
}

func initTofu(tempDir string) (string, error) {
	stdOut, stdErr, _, err := execute.Execute(
		"binaries/tofu",
		[]string{
			"-chdir=" + tempDir,
			"init",
			"-input=false",
			"-no-color"},
		map[string]string{})

	if err != nil {
		return "", errors.New("Failed to init: " + stdErr + " " + stdOut + " " + err.Error())
	}

	lockFile, err := os.ReadFile(filepath.Join(tempDir, ".terraform.lock.hcl"))

	if err != nil {
		return "", errors.New("Failed to get lock file: " + stdErr)
	}

	return base64.StdEncoding.EncodeToString(lockFile), nil
}

func generatePlan(tempDir string, token string, aud string, spaceId string) (string, string, error) {
	planFile := filepath.Join(tempDir, "tfplan")

	stdOut, stdErr, _, err := execute.Execute(
		"binaries/tofu",
		[]string{
			"-chdir=" + tempDir,
			"plan",
			"-no-color",
			"-out",
			planFile,
			"-var=octopus_space_id=" + spaceId},
		map[string]string{
			"OCTOPUS_ACCESS_TOKEN": token,
			"OCTOPUS_URL":          aud,
		})

	if err != nil {
		return "", "", errors.New("Failed to generate plan: " + stdErr + " " + stdOut + " " + err.Error())
	}

	plan, err := os.ReadFile(planFile)

	if err != nil {
		return "", "", err
	}

	return planFile, base64.StdEncoding.EncodeToString(plan), nil
}

func generatePlanJson(planFile string) (string, error) {
	planJsonStdOut, stdErr, _, err := execute.Execute(
		"binaries/tofu",
		[]string{
			"show",
			"-json",
			"-no-color",
			planFile},
		map[string]string{})

	if err != nil {
		return "", errors.New("Failed to generate plan json: " + stdErr + " " + planJsonStdOut + " " + err.Error())
	}

	return planJsonStdOut, nil
}

func generatePlanText(planFile string) (string, error) {
	planStdOut, stdErr, _, err := execute.Execute(
		"binaries/tofu",
		[]string{
			"show",
			"-no-color",
			planFile},
		map[string]string{})

	if err != nil {
		return "", errors.New("Failed to generate plan text: " + stdErr + " " + planStdOut + " " + err.Error())
	}

	return planStdOut, nil
}

func checkPlan(planJson string) error {
	tempDir, err := files.CreateTempDir()

	if err != nil {
		return err
	}

	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			zap.L().Error("Failed to remove temporary directory", zap.Error(err))
		}
	}()

	if err := os.WriteFile(filepath.Join(tempDir, "plan.json"), []byte(planJson), 0644); err != nil {
		return err
	}

	checkStdOut, _, exitCode, err := execute.Execute(
		"binaries/opa_linux_amd64",
		[]string{
			"exec",
			"--fail",
			"--decision",
			"terraform/analysis/allow",
			"--bundle",
			"policy/",
			filepath.Join(tempDir, "plan.json")},
		map[string]string{})

	if err != nil {
		return err
	}

	if exitCode != 0 {
		return fmt.Errorf("OPA check failed with exit code %d: %s", exitCode, checkStdOut)
	}

	zap.L().Info(checkStdOut)

	return nil
}

// createTerraformRcFile creates a .terraformrc file in the user's home directory
// The providers directory structure needs to be like:
// provider/registry.terraform.io/octopusdeploylabs/octopusdeploy/0.41.0/linux_amd64/terraform-provider-octopusdeploy_v0.41.0
func createTerraformRcFile() error {
	currentDir, err := os.Getwd()

	if err != nil {
		return fmt.Errorf("failed to determine current working directory: %w", err)
	}

	content := `provider_installation {
  filesystem_mirror {
    path    = "` + currentDir + `/provider"
    include = ["*/*/*"]
  }
  direct {}
}`

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to determine user home directory: %w", err)
	}

	rcFilePath := filepath.Join(homeDir, ".terraformrc")

	// Check if file already exists
	if err := backupRcFile(rcFilePath); err != nil {
		return err
	}

	// Write the new content
	if err := os.WriteFile(rcFilePath, []byte(content), 0600); err != nil {
		return fmt.Errorf("failed to write .terraformrc file: %w", err)
	}

	return nil
}

func backupRcFile(rcFilePath string) error {
	if _, err := os.Stat(rcFilePath); err == nil {
		// Back up existing file
		backupPath := rcFilePath + ".backup." + time.Now().Format("20060102150405")
		if err := os.Rename(rcFilePath, backupPath); err != nil {
			return fmt.Errorf("failed to backup existing .terraformrc file: %w", err)
		}
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("failed to check if .terraformrc file exists: %w", err)
	}

	return nil
}
