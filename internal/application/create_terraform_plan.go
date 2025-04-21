package application

import (
	"github.com/DataDog/jsonapi"
	"github.com/OctopusSolutionsEngineering/OctoAISpaceBuilder/internal/application/responses"
	"github.com/OctopusSolutionsEngineering/OctoAISpaceBuilder/internal/domain/handler"
	"github.com/OctopusSolutionsEngineering/OctoAISpaceBuilder/internal/domain/model"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

func CreateTerraformPlan(c *gin.Context) {

	body, err := io.ReadAll(c.Request.Body)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, responses.GenerateError("Failed to process request", err))
		return
	}

	var terraformInput model.TerraformPlan
	err = jsonapi.Unmarshal(body, &terraformInput)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, responses.GenerateError("Failed to process request", err))
		return
	}

	server, token, apiKey, err := getServerTokenApiKey(c)

	response, err := handler.CreateTerraformPlan(server, token, apiKey, terraformInput)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, responses.GenerateError("Failed to process request", err))
		return
	}

	responseJSON, err := jsonapi.Marshal(response)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, responses.GenerateError("Failed to process request", err))
		return
	}

	c.Header("Content-Type", "application/json")
	c.String(http.StatusCreated, string(responseJSON))
}
