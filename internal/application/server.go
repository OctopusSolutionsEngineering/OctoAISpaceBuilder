package application

import (
	"github.com/OctopusSolutionsEngineering/OctoAISpaceBuilder/internal/application/middleware"
	"github.com/OctopusSolutionsEngineering/OctoAISpaceBuilder/internal/domain/environment"
	"github.com/gin-gonic/gin"
)

func StartServer() error {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.POST("/api/terraformplan", middleware.JwtCheckMiddleware(environment.DisableValidation()), middleware.MakeExecutable, CreateTerraformPlan)
	router.POST("/api/terraformapply", middleware.JwtCheckMiddleware(environment.DisableValidation()), middleware.MakeExecutable, CreateTerraformApply)
	router.GET("/api/health", Health)

	return router.Run("localhost:" + environment.GetPort())
}
