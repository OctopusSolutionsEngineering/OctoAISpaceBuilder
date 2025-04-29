package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/OctopusSolutionsEngineering/OctoAISpaceBuilder/internal/domain/customerrors"
	"github.com/OctopusSolutionsEngineering/OctoAISpaceBuilder/internal/domain/environment"
	"github.com/OctopusSolutionsEngineering/OctoAISpaceBuilder/internal/domain/execute"
	"github.com/OctopusSolutionsEngineering/OctoAISpaceBuilder/internal/domain/files"
	"github.com/OctopusSolutionsEngineering/OctoAISpaceBuilder/internal/domain/logging"
	"github.com/OctopusSolutionsEngineering/OctoAISpaceBuilder/internal/domain/model"
	"github.com/OctopusSolutionsEngineering/OctoAISpaceBuilder/internal/domain/sha"
	"github.com/OctopusSolutionsEngineering/OctoAISpaceBuilder/internal/domain/terraform"
	"github.com/OctopusSolutionsEngineering/OctoAISpaceBuilder/internal/domain/validation"
	"github.com/OctopusSolutionsEngineering/OctoAISpaceBuilder/internal/infrastructure"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"os"
	"path/filepath"
	"time"
)

func CreateTerraformPlan(server string, token string, apiKey string, terraformInput model.TerraformPlan) (*model.TerraformPlan, error) {

	logging.LogEnhanced(terraformInput.Configuration, server)

	if err := validation.ValidateTerraformPlanRequest(terraformInput); err != nil {
		return nil, err
	}

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

	if err := terraform.WriteOverrides(tempDir); err != nil {
		return nil, err
	}

	if !environment.GetDisableTerraformCliConfig() {
		if err := terraform.CreateTerraformRcFile(); err != nil {
			return nil, err
		}
	}

	lockFile, err := initTofu(tempDir)

	if err != nil {
		return nil, err
	}

	planFile, planBinary, err := generatePlan(tempDir, token, apiKey, server, terraformInput.SpaceId)

	if err != nil {
		return nil, err
	}

	response := model.TerraformPlan{
		ID:      uuid.New().String(),
		Created: time.Now(),
		Server:  sha.GetSha256Hash(server),
		SpaceId: terraformInput.SpaceId,
	}

	if err := infrastructure.CreatePlanAzureStorageBlob(response.ID, planBinary, lockFile, []byte(terraformInput.Configuration)); err != nil {
		return nil, err
	}

	planJson, err := generatePlanJson(tempDir, planFile)

	if err != nil {
		return nil, err
	}

	if err := checkPlan(planJson, server); err != nil {
		return nil, err
	}

	planText, err := generatePlanText(tempDir, planFile)

	if err != nil {
		return nil, err
	}

	response.PlanText = &planText

	return &response, nil
}

func initTofu(tempDir string) ([]byte, error) {
	zap.L().Info("Init tofu")

	stdOut, stdErr, _, err := execute.Execute(
		environment.GetTofuExecutable(),
		[]string{
			"-chdir=" + tempDir,
			"init",
			"-input=false",
			"-no-color"},
		map[string]string{
			"TF_INPUT": "0",
		})

	if err != nil {
		return nil, errors.New("Failed to init: " + stdErr + " " + stdOut + " " + err.Error())
	}

	lockFile, err := os.ReadFile(filepath.Join(tempDir, ".terraform.lock.hcl"))

	if err != nil {
		return nil, errors.New("Failed to get lock file: " + stdErr)
	}

	return lockFile, nil
}

func generatePlan(tempDir string, token string, apiKey string, aud string, spaceId string) (string, []byte, error) {
	zap.L().Info("Generating plan for " + aud)

	planFile := filepath.Join(tempDir, "tfplan")

	stdOut, stdErr, _, err := execute.Execute(
		environment.GetTofuExecutable(),
		[]string{
			"-chdir=" + tempDir,
			"plan",
			"-no-color",
			"-out",
			planFile,
			"-var=octopus_space_id=" + spaceId},
		map[string]string{
			"OCTOPUS_ACCESS_TOKEN":  token,
			"OCTOPUS_API_KEY":       apiKey,
			"OCTOPUS_URL":           aud,
			"TF_INPUT":              "0",
			"TF_VAR_octopus_apikey": "",
			"TF_VAR_octopus_server": "",
		})

	if err != nil {
		logging.LogEnhanced(stdErr, aud)
		logging.LogEnhanced(stdOut, aud)
		return "", nil, errors.New("Failed to generate plan: " + stdErr + " " + stdOut + " " + err.Error())
	}

	plan, err := os.ReadFile(planFile)

	if err != nil {
		return "", nil, err
	}

	return planFile, plan, nil
}

func generatePlanJson(tempDir string, planFile string) (string, error) {
	zap.L().Info("Generating plan JSON")

	planJsonStdOut, stdErr, _, err := execute.Execute(
		environment.GetTofuExecutable(),
		[]string{
			"-chdir=" + tempDir,
			"show",
			"-json",
			"-no-color",
			planFile},
		map[string]string{
			"TF_INPUT": "0",
		})

	if err != nil {
		return "", errors.New("Failed to generate plan json: " + stdErr + " " + planJsonStdOut + " " + err.Error())
	}

	return planJsonStdOut, nil
}

func generatePlanText(tempDir string, planFile string) (string, error) {
	zap.L().Info("Generating plan text")

	planStdOut, stdErr, _, err := execute.Execute(
		environment.GetTofuExecutable(),
		[]string{
			"-chdir=" + tempDir,
			"show",
			"-no-color",
			planFile},
		map[string]string{
			"TF_INPUT": "0",
		})

	if err != nil {
		return "", errors.New("Failed to generate plan text: " + stdErr + " " + planStdOut + " " + err.Error())
	}

	return planStdOut, nil
}

func checkPlan(planJson string, server string) error {
	zap.L().Info("Checking plan with OPA")

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
		environment.GetOpaExecutable(),
		[]string{
			"exec",
			"--fail",
			"--decision",
			"terraform/analysis/allow",
			"--bundle",
			environment.GetOpaPolicyPath(),
			filepath.Join(tempDir, "plan.json")},
		map[string]string{})

	if err != nil {
		return err
	}

	if exitCode != 0 {
		return fmt.Errorf("OPA check failed with exit code %d: %s", exitCode, checkStdOut)
	}

	// Parse the OPA JSON output
	var opaResponse model.OpaResult
	if err := json.Unmarshal([]byte(checkStdOut), &opaResponse); err != nil {
		return fmt.Errorf("failed to parse OPA response: %w", err)
	}

	// Check the result from the parsed JSON
	for _, result := range opaResponse.Result {
		if !result.Result {
			logging.LogEnhanced(checkStdOut, server)
			return customerrors.OpaValidationFailed{
				ExitCode:   exitCode,
				DecisionID: result.DecisionID,
				Path:       result.Path,
			}
		}
	}

	zap.L().Info(checkStdOut)

	return nil
}
