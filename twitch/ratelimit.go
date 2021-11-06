package twitch

import (
	"github.com/nicklaw5/helix"
	"github.com/rs/zerolog/log"
	"time"
)

// ratelimitCallback ...
// https://github.com/nicklaw5/helix/blob/main/docs/README.md#request-rate-limiting
func ratelimitCallback(res *helix.Response) error {
	if res.GetRateLimitRemaining() > 0 {
		return nil
	}

	reset := int64(res.GetRateLimitReset())
	t := time.Now().Unix()

	if t < reset {
		untilReset := time.Duration(reset-t) * time.Second
		if untilReset > 0 {
			log.Info().Msgf("Hit ratelimit, waiting %s", untilReset.String())
			time.Sleep(untilReset)
		}
	}

	return nil
}
