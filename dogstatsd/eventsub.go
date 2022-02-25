package dogstatsd

import (
	"github.com/rs/zerolog/log"
	"github.com/streamcord/helix/v2"
)

// LogEventSubSubscriptionCount - Submit the amount of EventSub subscriptions used and allowed by Twitch.
func LogEventSubSubscriptionCount(res *helix.EventSubSubscriptionsResponse, service string) {
	if res == nil {
		return
	}

	log.Debug().
		Str("service", service).
		Msgf("%d of %d eventsub subscriptions used", res.Data.TotalCost, res.Data.MaxTotalCost)

	err := GlobalClient.Gauge(
		MetricTwitchEventSubSubscriptionsUsed,
		float64(res.Data.TotalCost),
		[]string{
			"service:" + service,
		},
		GlobalRate,
	)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to submit DogStatsD packet")
	}

	err = GlobalClient.Gauge(
		MetricTwitchEventSubSubscriptionsTotal,
		float64(res.Data.MaxTotalCost),
		[]string{
			"service:" + service,
		},
		GlobalRate,
	)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to submit DogStatsD packet")
	}
}
