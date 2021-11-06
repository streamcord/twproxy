package twitch

import (
	"github.com/nicklaw5/helix"
	"github.com/rs/zerolog/log"
	"time"
)

// EnsureAppAccessToken requests an App Access Token and creates a background task to automatically refresh the token
// before it expires.
func EnsureAppAccessToken(svcName string, c *helix.Client) {
	// Request and set an app access token
	creds := initialRequestAndSet(svcName, c)

	// Start the oauth token refresh service
	go func() {
		for {
			// Convert expiry time to time.Duration, and add a safety of 60 seconds to allow for token refresh
			expiresIn := time.Duration(creds.ExpiresIn-60) * time.Second
			log.Info().
				Str("service", svcName).
				Msgf(
					"Access token will expire in %s (~%v days)",
					expiresIn.String(),
					int(expiresIn.Hours())/24)

			time.Sleep(expiresIn)

			creds = refreshToken(svcName, c, creds.RefreshToken)
		}
	}()
}

func initialRequestAndSet(svcName string, c *helix.Client) helix.AccessCredentials {
	res, err := c.RequestAppAccessToken([]string{})
	if err != nil {
		log.Fatal().
			Str("service", svcName).
			Err(err).
			Msg("Failed to request app access token")
	} else if res.StatusCode != 200 {
		log.Fatal().
			Str("service", svcName).
			Int("status", res.StatusCode).
			Interface("data", res).
			Msg("Failed to request app access token")
	}

	c.SetAppAccessToken(res.Data.AccessToken)
	log.Debug().
		Str("service", svcName).
		Msgf("Set app access token, expires in %d seconds", res.Data.ExpiresIn)

	return res.Data
}

func refreshToken(svcName string, c *helix.Client, refreshToken string) helix.AccessCredentials {
	res, err := c.RefreshUserAccessToken(refreshToken)
	if err != nil {
		log.Fatal().
			Str("service", svcName).
			Err(err).
			Msg("Failed to refresh app access token")
	} else if res.StatusCode != 200 {
		log.Fatal().
			Str("service", svcName).
			Int("status", res.StatusCode).
			Interface("data", res).
			Msg("Failed to refresh app access token")
	}

	c.SetAppAccessToken(res.Data.AccessToken)
	log.Debug().
		Str("service", svcName).
		Msgf("Refreshed app access token, expires in %d seconds", res.Data.ExpiresIn)

	return res.Data
}
