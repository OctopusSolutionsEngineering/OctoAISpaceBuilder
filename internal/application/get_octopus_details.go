package application

import (
	"github.com/OctopusSolutionsEngineering/OctoAISpaceBuilder/internal/application/responses"
	"github.com/OctopusSolutionsEngineering/OctoAISpaceBuilder/internal/domain/jwt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func getServerTokenApiKey(c *gin.Context) (string, string, string, error) {
	token := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")

	if token != "" {
		server, err := jwt.GetJwtAud(token)

		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, responses.GenerateError("Failed to process request", err))
			return "", "", "", err
		}

		return server, token, "", nil
	}

	apiKey := c.GetHeader("X-Octopus-ApiKey")
	server := c.GetHeader("X-Octopus-Url")

	return server, "", apiKey, nil
}
