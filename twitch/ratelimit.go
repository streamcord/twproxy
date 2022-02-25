package twitch

import (
	"github.com/rs/zerolog/log"
	"github.com/streamcord/helix/v2"
	"time"
)

func (c *Client) genRatelimitFunc() func(res *helix.Response) error {
	return func(res *helix.Response) error {
		if res.GetRateLimitRemaining() > 0 {
			return nil
		}

		reset := int64(res.GetRateLimitReset())
		t := time.Now().Unix()

		if t < reset {
			untilReset := time.Duration(reset-t) * time.Second
			if untilReset > 0 {
				log.Info().
					Str("service", c.Service).
					Msgf("Hit ratelimit, waiting %s", untilReset.String())
				time.Sleep(untilReset)
			}
		}

		return nil
	}
}
