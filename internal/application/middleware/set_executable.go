package middleware

import (
	"github.com/OctopusSolutionsEngineering/OctoAISpaceBuilder/internal/application/responses"
	"github.com/OctopusSolutionsEngineering/OctoAISpaceBuilder/internal/domain/execute"
	"github.com/gin-gonic/gin"
	"net/http"
)

// MakeExecutable is a middleware function that sets the executable permissions for the specified files.
func MakeExecutable(c *gin.Context) {
	if err := execute.MakeAllExecutable("binaries"); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, responses.GenerateError("Failed to process request", err))
		c.Abort()
		return
	}

	if err := execute.MakeAllExecutable("provider"); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, responses.GenerateError("Failed to process request", err))
		c.Abort()
		return
	}
}
