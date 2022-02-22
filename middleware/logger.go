package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"time"
)

// LoggerMiddleware ...
func LoggerMiddleware(c *gin.Context) {
	start := time.Now()
	c.Next()
	d := time.Now().Sub(start)

	log.Info().
		Int("status", c.Writer.Status()).
		Str("method", c.Request.Method).
		Str("path", c.Request.URL.Path).
		Str("latency", d.String()).
		Str("remote_addr", c.ClientIP()).
		Msg("")
}
