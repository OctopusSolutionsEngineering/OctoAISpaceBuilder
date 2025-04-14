package application

import (
	"github.com/OctopusSolutionsEngineering/OctoAISpaceBuilder/internal/application/middleware"
	"github.com/OctopusSolutionsEngineering/OctoAISpaceBuilder/internal/domain/environment"
	"github.com/gin-gonic/gin"
)

func StartServer() error {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.POST("/api/terraformplan", middleware.JwtCheckMiddleware(environment.DisableValidation()), CreateTerraformPlan)
	router.POST("/api/terraformapply", middleware.JwtCheckMiddleware(environment.DisableValidation()), CreateTerraformApply)

	return router.Run("localhost:" + environment.GetPort())
}
