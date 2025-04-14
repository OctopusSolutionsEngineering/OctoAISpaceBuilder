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

func CreateTerraformApply(c *gin.Context) {

	token := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")

	body, err := io.ReadAll(c.Request.Body)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, responses.GenerateError("Failed to process request", err))
		return
	}

	var terraform model.TerraformApply
	err = jsonapi.Unmarshal(body, &terraform)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, responses.GenerateError("Failed to process request", err))
		return
	}

	result, err := handler.CreateTerraformApply(token, terraform)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, responses.GenerateError("Failed to process request", err))
		return
	}

	c.String(http.StatusCreated, result)
}
