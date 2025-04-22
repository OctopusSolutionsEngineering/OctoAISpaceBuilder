package logging

import (
	"github.com/OctopusSolutionsEngineering/OctoAISpaceBuilder/internal/domain/environment"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

func IsEnhancedLoggingEnabled(server string) bool {
	servers := environment.GetEnhancedLoggingInstances()
	return lo.Contains(servers, server)
}

// LogEnhanced will log a message only if the server that made the request is specifically configured
// for enhanced logging.
func LogEnhanced(log string, server string) {
	enhancedLogging := IsEnhancedLoggingEnabled(server)

	if enhancedLogging {
		zap.L().Info(log)
	}
}
