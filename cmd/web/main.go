package main

import (
	"fmt"
	"github.com/OctopusSolutionsEngineering/OctoAISpaceBuilder/internal/application"
	"github.com/OctopusSolutionsEngineering/OctoAISpaceBuilder/internal/domain/environment"
	"github.com/OctopusSolutionsEngineering/OctoAISpaceBuilder/internal/domain/logging"
	"go.uber.org/zap"
	"os"
)

func main() {
	logging.ConfigureZapLogger()

	cwd, err := os.Getwd()
	if err != nil {
		zap.L().Error(err.Error())
		return
	}

	zap.L().Info("Current working directory: " + cwd)
	zap.L().Info("Disable validation: " + fmt.Sprint(environment.DisableValidation()))
	zap.L().Info("OPA executable: " + fmt.Sprint(environment.GetOpaExecutable()))
	zap.L().Info("OPA policy path: " + fmt.Sprint(environment.GetOpaPolicyPath()))
	zap.L().Info("Tofu executable: " + fmt.Sprint(environment.GetTofuExecutable()))
	zap.L().Info("Disable setting binary execution flag: " + fmt.Sprint(environment.DisableMakeBinariesExecutable()))
	zap.L().Info("Enhanced logging instance: " + fmt.Sprint(environment.GetEnhancedLoggingInstances()))
	zap.L().Info("Port: " + fmt.Sprint(environment.GetPort()))

	if err := application.StartServer(); err != nil {
		zap.L().Error(err.Error())
	}
}
