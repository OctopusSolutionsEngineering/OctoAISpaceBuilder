package middleware

import (
	"github.com/gin-gonic/gin"
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

	c.IndentedJSON(http.StatusUnauthorized, "No authorization token or API key provided")
	c.Abort()
	return
}
