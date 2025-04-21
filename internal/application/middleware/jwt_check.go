package middleware

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusSolutionsEngineering/OctoAISpaceBuilder/internal/application/responses"
	"github.com/OctopusSolutionsEngineering/OctoAISpaceBuilder/internal/domain/jwt"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"strings"
)

func JwtCheckMiddleware(skipValidation bool) gin.HandlerFunc {
	return func(c *gin.Context) {

		// At the end of the day, this service is essentially unauthenticated.
		// We accept any user with a valid JWT token that appears to authenticate with an Octopus Deploy instance.
		token := strings.TrimSpace(strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer "))

		if token == "" {
			// If the token is empty, we don't need to do anything
			c.Next()
			return
		}

		aud, err := jwt.ValidateJWT(token, skipValidation)

		if err != nil {
			c.IndentedJSON(http.StatusUnauthorized, responses.GenerateError("Failed to validate token in middleware", err))
			c.Abort()
			return
		}

		apiURL, err := url.Parse(aud)
		if err != nil {
			c.IndentedJSON(http.StatusUnauthorized, responses.GenerateError("Failed to process request", err))
			c.Abort()
			return
		}

		// Don't try to validate the token with an API call
		if skipValidation {
			c.Next()
			return
		}

		// Use the token to look up the user. This is not foolproof - you could supply any valid JWT token
		// with an audience claim that points to a server that responds to this API request.
		// We can't prove that anyone submitting feedback is a genuine Octopus user.
		// But we do effectively prove that you own a DNS name, which is almost as good.
		// Since we store the audience in the feedback items, we can filter out bad requests later.
		// It also raises the bar for anyone looking to abuse the API, as you would need to generate valid JWTs,
		// host a JWKS server, and host a server that responds the API request.
		octopusClient, err := client.NewClientWithAccessToken(nil, apiURL, token, "")
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, responses.GenerateError("Failed to process request", err))
			c.Abort()
			return
		}

		if _, err := octopusClient.Users.GetMe(); err != nil {
			c.IndentedJSON(http.StatusUnauthorized, responses.GenerateError("Failed to process request", err))
			c.Abort()
			return
		}

		// normal request, and the execution chain is called down
		c.Next()
	}
}
