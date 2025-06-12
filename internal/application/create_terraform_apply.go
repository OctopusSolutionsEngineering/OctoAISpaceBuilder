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

	response, err := handler.CreateTerraformApply(server, token, apiKey, terraform)

	if err != nil {
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
