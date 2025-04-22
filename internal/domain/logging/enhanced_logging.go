package logging

import (
	"github.com/OctopusSolutionsEngineering/OctoAISpaceBuilder/internal/domain/environment"
	"github.com/samber/lo"
)

func IsEnhancedLoggingEnabled(server string) bool {
	servers := environment.GetEnhancedLoggingInstances()
	return lo.Contains(servers, server)
}
