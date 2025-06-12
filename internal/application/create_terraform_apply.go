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

func CreateTerraformApply(c *gin.Context) {

	body, err := io.ReadAll(c.Request.Body)

	if err != nil {
		zap.L().Error("Failed to read request body", zap.Error(err))
		c.IndentedJSON(http.StatusBadRequest, responses.GenerateError("Failed to process request", err))
		return
	}

	var terraform model.TerraformApply
	err = jsonapi.Unmarshal(body, &terraform)

	if err != nil {
		zap.L().Error("Failed to unmarshal JSON API body", zap.Error(err))
		c.IndentedJSON(http.StatusBadRequest, responses.GenerateError("Failed to process request", err))
		return
	}

	server, token, apiKey, err := getServerTokenApiKey(c)

	if err != nil {
		zap.L().Error("Failed to get the Octopus details", zap.Error(err))
		c.IndentedJSON(http.StatusBadRequest, responses.GenerateError("Failed to process request", err))
		return
	}

	var response *model.TerraformApply = nil
	var applyError error = nil

	// There are cases where the terraform apply might fail by the resources are created.
	// Errors like "unable to locate variable for owner ID Projects-2" are an example where
	// applying the plan a second time succeeds.
	// This is most likely due to a bug in the Terraform provider.
	// But, there is no harm in retrying the apply operation a second time.
	for i := 0; i < 2; i++ {
		response, applyError = handler.CreateTerraformApply(server, token, apiKey, terraform)
		if applyError == nil {
			break
		}
	}

	if applyError != nil {
		zap.L().Error("Failed to perform Terraform apply", zap.Error(err))
		c.IndentedJSON(http.StatusInternalServerError, responses.GenerateError("Failed to process request", err))
		return
	}

	responseJSON, err := jsonapi.Marshal(response)

	if err != nil {
		zap.L().Error("Failed to marshal JSON API response", zap.Error(err))
		c.IndentedJSON(http.StatusInternalServerError, responses.GenerateError("Failed to process request", err))
		return
	}

	c.Header("Content-Type", "application/vnd.api+json")
	c.String(http.StatusCreated, string(responseJSON))
}
