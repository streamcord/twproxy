package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"twproxy/config"
	"twproxy/twitch"
)

// AuthMiddleware - Handles basic authentication.
// Clients must provide two extra headers when sending requests: `Client-ID` and `Authorization`.
// These headers will be used to determine which helix.Client to use, and if the client is authorized to access the proxy.
func AuthMiddleware(c *gin.Context) {
	svcName := c.Request.Header.Get("Client-ID")
	if svcName == "" {
		// Default to "default" if no header was specified
		svcName = "default"
	}

	svc, ok := config.GlobalConfig.Services[svcName]
	if !ok {
		c.AbortWithStatusJSON(401, gin.H{
			"error":   "Unauthorized",
			"message": "Unknown service name '" + svcName + "'",
			"status":  401,
		})
		return
	}

	svcAuth := c.Request.Header.Get("Authorization")
	if svcAuth != "Bearer "+svc.Auth {
		log.Warn().Str("service", svcName).Str("auth", svcAuth).Msg("Invalid authorization attempt")
		c.AbortWithStatusJSON(403, gin.H{
			"error":   "Forbidden",
			"message": "Invalid auth for service '" + svcName + "'",
			"status":  403,
		})
		return
	}

	c.Set("helix", twitch.GlobalClients[svcName])
}
