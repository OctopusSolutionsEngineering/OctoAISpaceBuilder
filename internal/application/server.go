package application

import (
	"github.com/OctopusSolutionsEngineering/OctoAISpaceBuilder/internal/application/middleware"
	"github.com/OctopusSolutionsEngineering/OctoAISpaceBuilder/internal/domain/environment"
	"github.com/gin-gonic/gin"
)

func StartServer() error {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.POST("/api/terraform", middleware.JwtCheckMiddleware(environment.DisableValidation()), CreateTerraform)

	return router.Run("localhost:" + environment.GetPort())
}
