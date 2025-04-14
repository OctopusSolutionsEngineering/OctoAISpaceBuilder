package application

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/DataDog/jsonapi"
	"github.com/OctopusSolutionsEngineering/OctoAISpaceBuilder/internal/application/responses"
	"github.com/OctopusSolutionsEngineering/OctoAISpaceBuilder/internal/domain/execute"
	"github.com/OctopusSolutionsEngineering/OctoAISpaceBuilder/internal/domain/files"
	"github.com/OctopusSolutionsEngineering/OctoAISpaceBuilder/internal/domain/jwt"
	"github.com/OctopusSolutionsEngineering/OctoAISpaceBuilder/internal/domain/model"
	"github.com/OctopusSolutionsEngineering/OctoAISpaceBuilder/internal/domain/sha"
	"github.com/OctopusSolutionsEngineering/OctoAISpaceBuilder/internal/infrastructure"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func CreateTerraformPlan(c *gin.Context) {

	body, err := io.ReadAll(c.Request.Body)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, responses.GenerateError("Failed to process request", err))
		return
	}

	var terraform model.Terraform
	err = jsonapi.Unmarshal(body, &terraform)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, responses.GenerateError("Failed to process request", err))
		return
	}

	tempDir, err := files.CreateTempDir()

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, responses.GenerateError("Failed to process request", err))
		return
	}

	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			zap.L().Error("Failed to remove temporary directory", zap.Error(err))
		}
	}()

	if err := os.WriteFile(filepath.Join(tempDir, "terraform.tf"), []byte(terraform.Configuration), 0644); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, responses.GenerateError("Failed to process request", err))
		return
	}

	token := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")

	planFile, planBinary, err := generatePlan(tempDir, token)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, responses.GenerateError("Failed to process request", err))
		return
	}

	aud, err := jwt.GetJwtAud(token)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, responses.GenerateError("Failed to process request", err))
		return
	}

	terraformPlan := model.TerraformPlan{
		ID:               uuid.New().String(),
		PlanBinaryBase64: &planBinary,
		Created:          time.Now(),
		Server:           sha.GetSha256Hash(aud),
	}

	if err := infrastructure.CreateFeedbackAzureStorageTable(terraformPlan); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, responses.GenerateError("Failed to process request", err))
		return
	}

	planJson, err := generatePlanJson(planFile)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, responses.GenerateError("Failed to process request", err))
		return
	}

	if err := checkPlan(planJson); err != nil {
		c.IndentedJSON(http.StatusBadRequest, responses.GenerateError("Failed to process request", err))
		return
	}

	planText, err := generatePlanText(planFile)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, responses.GenerateError("Failed to process request", err))
		return
	}

	response := model.TerraformPlan{
		ID:       terraformPlan.ID,
		PlanText: &planText,
		Created:  terraformPlan.Created,
		Server:   terraformPlan.Server,
	}

	responseJSON, err := jsonapi.Marshal(response)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, responses.GenerateError("Failed to process request", err))
		return
	}

	c.String(http.StatusCreated, string(responseJSON))
}

func generatePlan(tempDir string, token string) (string, string, error) {
	planFile := filepath.Join(tempDir, "tfplan")

	_, stdErr, _, err := execute.Execute(
		"binaries/tofu",
		[]string{
			"-chdir=" + tempDir,
			"plan",
			"-no-color",
			"-out",
			planFile},
		map[string]string{
			"TF_VAR_access_token": token,
		})

	if err != nil {
		return "", "", errors.New("Failed to generate plan: " + stdErr)
	}

	plan, err := os.ReadFile(planFile)

	if err != nil {
		return "", "", err
	}

	return planFile, base64.StdEncoding.EncodeToString(plan), nil
}

func generatePlanJson(planFile string) (string, error) {
	planJsonStdOut, _, _, err := execute.Execute(
		"binaries/tofu",
		[]string{
			"show",
			"-json",
			"-no-color",
			planFile},
		map[string]string{})

	if err != nil {
		return "", nil
	}

	return planJsonStdOut, nil
}

func generatePlanText(planFile string) (string, error) {
	planStdOut, _, _, err := execute.Execute(
		"binaries/tofu",
		[]string{
			"show",
			"-no-color",
			planFile},
		map[string]string{})

	if err != nil {
		return "", nil
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
