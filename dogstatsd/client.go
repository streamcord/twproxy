package dogstatsd

import (
	"github.com/DataDog/datadog-go/v5/statsd"
	"github.com/rs/zerolog/log"
	"twproxy/config"
)

// GlobalClient ...
var GlobalClient *statsd.Client

// New ...
func New() {
	a := config.GlobalConfig.DogStatsD.Address

	var err error
	GlobalClient, err = statsd.New(
		a,
		statsd.WithNamespace(config.GlobalConfig.DogStatsD.Namespace),
		// statsd.WithExtendedClientSideAggregation(),
	)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create dogStatsD client")
	}

	log.Info().Msgf("Created DogStatsD client connected to %s", a)
}
