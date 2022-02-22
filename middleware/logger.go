package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"strconv"
	"time"
	"twproxy/dogstatsd"
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

	// Send request info to DogStatsD
	go func() {
		err := dogstatsd.GlobalClient.Timing(dogstatsd.MetricRequest, d, []string{
			"status:" + strconv.Itoa(c.Writer.Status()),
			"method:" + c.Request.Method,
			"path:" + c.Request.URL.Path,
		}, dogstatsd.GlobalRate)
		if err != nil {
			log.Warn().Err(err).Msg("Failed to submit DogStatsD packet")
		}
	}()
}
