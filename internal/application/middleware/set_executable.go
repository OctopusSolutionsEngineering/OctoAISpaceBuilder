package middleware

import (
	"github.com/OctopusSolutionsEngineering/OctoAISpaceBuilder/internal/application/responses"
	"github.com/OctopusSolutionsEngineering/OctoAISpaceBuilder/internal/domain/environment"
	"github.com/OctopusSolutionsEngineering/OctoAISpaceBuilder/internal/domain/execute"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"os"
	"path/filepath"
)

// MakeExecutable is a middleware function that sets the executable permissions for the specified files.
func MakeExecutable(c *gin.Context) {
	if environment.DisableMakeBinariesExecutable() {
		return
	}

	cwd, err := os.Getwd()
	if err != nil {
		zap.L().Error("Failed to get the current working directory", zap.Error(err))
		c.IndentedJSON(http.StatusInternalServerError, responses.GenerateError("Failed to process request", err))
		c.Abort()
		return
	}

	if err := execute.MakeAllExecutable(filepath.Join(cwd, "binaries")); err != nil {
		zap.L().Error("Failed to make binaries executable", zap.Error(err))
		c.IndentedJSON(http.StatusInternalServerError, responses.GenerateError("Failed to process request", err))
		c.Abort()
		return
	}

	if err := execute.MakeAllExecutable(filepath.Join(cwd, "provider")); err != nil {
		zap.L().Error("Failed to make provider executable", zap.Error(err))
		c.IndentedJSON(http.StatusInternalServerError, responses.GenerateError("Failed to process request", err))
		c.Abort()
		return
	}
}
