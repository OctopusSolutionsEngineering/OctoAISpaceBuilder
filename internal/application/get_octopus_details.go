package application

import (
	"github.com/gin-gonic/gin"
	"strings"
)

func getServerTokenApiKey(c *gin.Context) (string, string, string, error) {
	token := strings.TrimSpace(strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer "))
	server := c.GetHeader("X-Octopus-Url")
	apiKey := c.GetHeader("X-Octopus-ApiKey")

	// Tokens take precedence over API keys
	if token != "" {
		return server, token, "", nil
	}

	return server, "", apiKey, nil
}
