package middleware

import (
	"net/http"
	"path/filepath"

	"github.com/OctopusSolutionsEngineering/OctoAISpaceBuilder/internal/application/responses"
	"github.com/OctopusSolutionsEngineering/OctoAISpaceBuilder/internal/domain/environment"
	"github.com/OctopusSolutionsEngineering/OctoAISpaceBuilder/internal/domain/execute"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// MakeExecutable is a middleware function that sets the executable permissions for the specified files.
func MakeExecutable(c *gin.Context) {
	if environment.DisableMakeBinariesExecutable() {
		return
	}

	installDir, err := environment.GetInstallationDirectory()
	if err != nil {
		zap.L().Error("Failed to get the installation directory", zap.Error(err))
		c.IndentedJSON(http.StatusInternalServerError, responses.GenerateError("Failed to process request", err))
		c.Abort()
		return
	}

	if err := execute.MakeAllExecutable(filepath.Join(installDir, "binaries")); err != nil {
		zap.L().Error("Failed to make binaries executable under "+installDir, zap.Error(err))
		c.IndentedJSON(http.StatusInternalServerError, responses.GenerateError("Failed to process request", err))
		c.Abort()
		return
	}

	if err := execute.MakeAllExecutable(filepath.Join(installDir, "provider")); err != nil {
		zap.L().Error("Failed to make provider executable under "+installDir, zap.Error(err))
		c.IndentedJSON(http.StatusInternalServerError, responses.GenerateError("Failed to process request", err))
		c.Abort()
		return
	}
}
