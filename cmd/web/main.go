package main

import (
	"github.com/OctopusSolutionsEngineering/OctoAISpaceBuilder/internal/application"
	"go.uber.org/zap"
	"os"
)

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		zap.L().Error(err.Error())
		return
	}

	zap.L().Info("Current working directory: " + cwd)

	if err := application.StartServer(); err != nil {
		zap.L().Error(err.Error())
	}
}
