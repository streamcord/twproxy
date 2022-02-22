package dogstatsd

import (
	"github.com/nicklaw5/helix"
	"github.com/rs/zerolog/log"
	"strconv"
	"time"
)

// Twitch API routes
const (
	RouteGetClips        = "/helix/clips"
	RouteGetGames        = "/helix/games"
	RouteGetStreams      = "/helix/streams"
	RouteGetUsers        = "/helix/users"
	RouteGetUsersFollows = "/helix/users/follows"
	RouteGetOauthToken   = "/oauth/token"
)

// parseRLInfo parses ratelimit info from a request into a string.
// Format: remaining/total
func parseRLInfo(res helix.ResponseCommon) string {
	return strconv.Itoa(res.GetRateLimitRemaining()) + "/" + strconv.Itoa(res.GetRateLimit())
}

// LogTwitchRequest - Submit Twitch HTTP request info to datadog.
func LogTwitchRequest(route string, service string, res helix.ResponseCommon, err error, latency time.Duration) {
	if err != nil {
		log.Error().
			Err(err).
			Str("service", service).
			Msgf("Twitch request to %s failed", route)
	} else if res.StatusCode > 399 {
		log.Warn().
			Str("latency", latency.String()).
			Str("ratelimit", parseRLInfo(res)).
			Str("service", service).
			Int("status", res.StatusCode).
			Msgf("Twitch request to %s failed: %s (%d): %s", route, res.Error, res.ErrorStatus, res.ErrorMessage)
	} else {
		log.Debug().
			Str("latency", latency.String()).
			Str("ratelimit", parseRLInfo(res)).
			Str("service", service).
			Int("status", res.StatusCode).
			Msgf("Executed Twitch request to %s", route)
	}

	// Submit request info to datadog
	t := []string{
		"route:" + route,
		"service:" + service,
		"status:" + strconv.Itoa(res.StatusCode),
	}
	ddErr := GlobalClient.Timing(MetricTwitchRequest, latency, t, GlobalRate)
	if ddErr != nil {
		log.Warn().Err(err).Msg("Failed to submit DogStatsD packet")
	}

	// Submit ratelimit info to datadog only if the request succeeded, i.e. we actually got a response from Twitch
	if err != nil {
		ddErr = GlobalClient.Gauge(
			MetricTwitchRatelimitRemaining,
			float64(res.GetRateLimitRemaining()),
			[]string{},
			GlobalRate,
		)
		if ddErr != nil {
			log.Warn().Err(err).Msg("Failed to submit DogStatsD packet")
		}
	}
}
