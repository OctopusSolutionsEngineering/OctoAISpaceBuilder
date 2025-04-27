package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

func AuthCheck(c *gin.Context) {
	token := c.GetHeader("Authorization")
	apiKey := c.GetHeader("X-Octopus-ApiKey")
	server := c.GetHeader("X-Octopus-Url")

	if token != "" {
		c.Next()
		return
	}

	if apiKey != "" && server != "" {
		c.Next()
		return
	}

	zap.L().Error("No authorization token or API key provided")
	c.IndentedJSON(http.StatusUnauthorized, "No authorization token or API key provided")
	c.Abort()
	return
}
