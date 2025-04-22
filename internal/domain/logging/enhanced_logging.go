package logging

import (
	"github.com/OctopusSolutionsEngineering/OctoAISpaceBuilder/internal/domain/environment"
	"github.com/samber/lo"
	"go.uber.org/zap"
	"net/url"
	"strings"
)

func getHostname(server string) string {
	parsedURL, err := url.Parse(server)
	if err != nil {
		return ""
	}
	return parsedURL.Hostname()
}

func IsEnhancedLoggingEnabled(server string) bool {
	servers := environment.GetEnhancedLoggingInstances()
	return lo.Contains(servers, getHostname(server))
}

// LogEnhanced will log a message only if the server that made the request is specifically configured
// for enhanced logging.
func LogEnhanced(log string, server string) {
	enhancedLogging := IsEnhancedLoggingEnabled(server)

	if enhancedLogging {
		// Long logs were getting truncated in the console, so we split them into lines
		for _, line := range strings.Split(log, "\n") {
			zap.L().Info(line)
		}
	}
}
