package application

import (
	"github.com/OctopusSolutionsEngineering/OctoAISpaceBuilder/internal/application/middleware"
	"github.com/OctopusSolutionsEngineering/OctoAISpaceBuilder/internal/domain/environment"
	"github.com/gin-gonic/gin"
)

func StartServer() error {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.POST("/api/terraformplan", middleware.AuthCheck, middleware.JwtCheckMiddleware(environment.DisableValidation()), middleware.CopyToWritablePath, middleware.MakeExecutable, CreateTerraformPlan)
	router.POST("/api/terraformapply", middleware.AuthCheck, middleware.JwtCheckMiddleware(environment.DisableValidation()), middleware.CopyToWritablePath, middleware.MakeExecutable, CreateTerraformApply)
	router.POST("/api/terraformautoapply", middleware.AuthCheck, middleware.JwtCheckMiddleware(environment.DisableValidation()), middleware.CopyToWritablePath, middleware.MakeExecutable, CreateTerraformAutoApply)
	router.GET("/api/health", Health)
	router.GET("/", Health)
	router.Any("/cleanup", CleanupOldPlans)

	return router.Run("0.0.0.0:" + environment.GetPort())
}
