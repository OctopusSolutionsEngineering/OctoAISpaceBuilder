package handler

import (
	"github.com/OctopusSolutionsEngineering/OctoAISpaceBuilder/internal/domain/compress"
	"github.com/OctopusSolutionsEngineering/OctoAISpaceBuilder/internal/domain/execute"
	"github.com/OctopusSolutionsEngineering/OctoAISpaceBuilder/internal/domain/files"
	"github.com/OctopusSolutionsEngineering/OctoAISpaceBuilder/internal/domain/model"
	"github.com/OctopusSolutionsEngineering/OctoAISpaceBuilder/internal/domain/sha"
	"github.com/OctopusSolutionsEngineering/OctoAISpaceBuilder/internal/domain/terraform"
	"github.com/OctopusSolutionsEngineering/OctoAISpaceBuilder/internal/infrastructure"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"os"
	"path/filepath"
)

func CreateTerraformApply(server string, token string, apiKey string, terraformApply model.TerraformApply) (*model.TerraformApply, error) {

	terraformApply.Server = sha.GetSha256Hash(server)

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

	planContents, _, lockFile, configuration, err := infrastructure.ReadPlanAzureStorageTable(terraformApply)

	if err != nil {
		return nil, err
	}

	decoded, err := compress.DecompressStringToByteArray(planContents)

	if err != nil {
		return nil, err
	}

	decodedLockFile, err := compress.DecompressStringToByteArray(lockFile)

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

	if err := terraform.WriteOverrides(tempDir); err != nil {
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
			"OCTOPUS_ACCESS_TOKEN":  token,
			"OCTOPUS_API_KEY":       apiKey,
			"OCTOPUS_URL":           server,
			"TF_INPUT":              "0",
			"TF_VAR_octopus_apikey": "",
			"TF_VAR_octopus_server": "",
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
			"OCTOPUS_ACCESS_TOKEN":  token,
			"OCTOPUS_URL":           server,
			"OCTOPUS_API_KEY":       apiKey,
			"TF_INPUT":              "0",
			"TF_VAR_octopus_apikey": "",
			"TF_VAR_octopus_server": "",
		})

	if err != nil {
		return nil, err
	}

	if err := infrastructure.DeleteTerraformPlan(infrastructure.TableEntityId{
		RowKey:       terraformApply.PlanId,
		PartitionKey: terraformApply.Server,
	}); err != nil {
		// We're not going to fail here, but we'll log the error.
		// Any old plans will be cleaned up by the terraform plan cleanup job.
		zap.L().Error("Failed to delete terraform plan", zap.Error(err))
	}

	response := model.TerraformApply{
		ID:        uuid.New().String(),
		PlanId:    terraformApply.PlanId,
		Server:    terraformApply.Server,
		ApplyText: &stdout,
	}

	return &response, nil
}
