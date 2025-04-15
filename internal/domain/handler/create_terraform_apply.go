package handler

import (
	"encoding/base64"
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
)

func CreateTerraformApply(token string, terraformApply model.TerraformApply) (*model.TerraformApply, error) {

	aud, err := jwt.GetJwtAud(token)

	if err != nil {
		return nil, err
	}

	terraformApply.Server = sha.GetSha256Hash(aud)

	tempDir, err := files.CreateTempDir()

	if err != nil {
		return nil, err
	}

	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			zap.L().Error("Failed to remove temporary directory", zap.Error(err))
		}
	}()

	planFile := filepath.Join(tempDir, "tfplan")
	lockFileName := filepath.Join(tempDir, ".terraformApply.lock.hcl")
	configurationFileName := filepath.Join(tempDir, "terraformApply.tf")

	planContents, _, lockFile, configuration, err := infrastructure.ReadFeedbackAzureStorageTable(terraformApply)

	if err != nil {
		return nil, err
	}

	decoded, err := base64.StdEncoding.DecodeString(planContents)

	if err != nil {
		return nil, err
	}

	decodedLockFile, err := base64.StdEncoding.DecodeString(lockFile)

	if err != nil {
		return nil, err
	}

	if err := os.WriteFile(planFile, decoded, 0644); err != nil {
		return nil, err
	}

	if err := os.WriteFile(lockFileName, decodedLockFile, 0644); err != nil {
		return nil, err
	}

	if err := os.WriteFile(lockFileName, decodedLockFile, 0644); err != nil {
		return nil, err
	}

	if err := os.WriteFile(configurationFileName, []byte(configuration), 0644); err != nil {
		return nil, err
	}

	if err := terraform.WriteBackendOverride(tempDir); err != nil {
		return nil, err
	}

	_, _, _, err = execute.Execute(
		"binaries/tofu",
		[]string{
			"-chdir=" + tempDir,
			"init",
			"-input=false",
			"-no-color"},
		map[string]string{
			"OCTOPUS_ACCESS_TOKEN": token,
			"OCTOPUS_URL":          aud,
		})

	if err != nil {
		return nil, err
	}

	stdout, _, _, err := execute.Execute(
		"binaries/tofu",
		[]string{
			"-chdir=" + tempDir,
			"apply",
			"-auto-approve",
			"-input=false",
			"-no-color",
			planFile},
		map[string]string{
			"OCTOPUS_ACCESS_TOKEN": token,
			"OCTOPUS_URL":          aud,
		})

	if err != nil {
		return nil, err
	}

	response := model.TerraformApply{
		ID:        uuid.New().String(),
		PlanId:    terraformApply.PlanId,
		Server:    terraformApply.Server,
		ApplyText: &stdout,
	}

	return &response, nil
}
