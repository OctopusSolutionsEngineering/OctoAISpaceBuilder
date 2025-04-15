package application

import (
	"github.com/DataDog/jsonapi"
	"github.com/OctopusSolutionsEngineering/OctoAISpaceBuilder/internal/application/responses"
	"github.com/OctopusSolutionsEngineering/OctoAISpaceBuilder/internal/domain/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

func Health(c *gin.Context) {
	responseJSON, err := jsonapi.Marshal(model.Health{
		ID:     uuid.New().String(),
		Status: "OK",
	})

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, responses.GenerateError("Failed to process request", err))
		return
	}

	c.String(200, string(responseJSON))
}
