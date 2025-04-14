package application

import (
	"github.com/DataDog/jsonapi"
	"github.com/OctopusSolutionsEngineering/OctoAISpaceBuilder/internal/application/responses"
	"github.com/OctopusSolutionsEngineering/OctoAISpaceBuilder/internal/domain/handler"
	"github.com/OctopusSolutionsEngineering/OctoAISpaceBuilder/internal/domain/model"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strings"
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

	token := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")

	response, err := handler.CreateTerraformPlan(token, terraformInput)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, responses.GenerateError("Failed to process request", err))
		return
	}

	c.String(http.StatusCreated, response)
}
