package application

import (
	"github.com/OctopusSolutionsEngineering/OctoAISpaceBuilder/internal/application/middleware"
	"github.com/OctopusSolutionsEngineering/OctoAISpaceBuilder/internal/domain/environment"
	"github.com/OctopusSolutionsEngineering/OctoAISpaceBuilder/internal/domain/logging"
	"github.com/gin-gonic/gin"
)

func StartServer() error {
	logging.ConfigureZapLogger()
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.POST("/api/terraformplan", middleware.AuthCheck, middleware.JwtCheckMiddleware(environment.DisableValidation()), middleware.MakeExecutable, CreateTerraformPlan)
	router.POST("/api/terraformapply", middleware.AuthCheck, middleware.JwtCheckMiddleware(environment.DisableValidation()), middleware.MakeExecutable, CreateTerraformApply)
	router.GET("/api/health", Health)
	router.GET("/", Health)
	router.Any("/cleanup", CleanupOldPlans)

	return router.Run("localhost:" + environment.GetPort())
}
