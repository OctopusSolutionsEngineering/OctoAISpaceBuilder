package application

import (
	"github.com/DataDog/jsonapi"
	"github.com/OctopusSolutionsEngineering/OctoAISpaceBuilder/internal/application/responses"
	"github.com/OctopusSolutionsEngineering/OctoAISpaceBuilder/internal/domain/handler"
	"github.com/OctopusSolutionsEngineering/OctoAISpaceBuilder/internal/domain/model"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io"
	"net/http"
)

func CreateTerraformAutoApply(c *gin.Context) {

	body, err := io.ReadAll(c.Request.Body)

	if err != nil {
		zap.L().Error("Failed to read request body", zap.Error(err))
		c.IndentedJSON(http.StatusBadRequest, responses.GenerateError("Failed to process request", err))
		return
	}

	var terraformInput model.TerraformPlan
	err = jsonapi.Unmarshal(body, &terraformInput)

	if err != nil {
		zap.L().Error("Failed to unmarshal body as JSON API", zap.Error(err))
		c.IndentedJSON(http.StatusBadRequest, responses.GenerateError("Failed to process request", err))
		return
	}

	server, token, apiKey, err := getServerTokenApiKey(c)

	response, err := handler.CreateTerraformPlan(server, token, apiKey, terraformInput)

	if err != nil {
		zap.L().Error("Failed to create Terraform plan", zap.Error(err))
		c.IndentedJSON(http.StatusInternalServerError, responses.GenerateError("Failed to process request", err))
		return
	}

	terraformApply := model.TerraformApply{
		PlanId: response.ID,
	}

	applyResponse, err := handler.CreateTerraformApply(server, token, apiKey, terraformApply)

	if err != nil {
		zap.L().Error("Failed to perform Terraform apply", zap.Error(err))
		c.IndentedJSON(http.StatusInternalServerError, responses.GenerateError("Failed to process request", err))
		return
	}

	responseJSON, err := jsonapi.Marshal(applyResponse)

	if err != nil {
		zap.L().Error("Failed to marshal JSON API response", zap.Error(err))
		c.IndentedJSON(http.StatusInternalServerError, responses.GenerateError("Failed to process request", err))
		return
	}

	c.Header("Content-Type", "application/vnd.api+json")
	c.String(http.StatusCreated, string(responseJSON))
}
