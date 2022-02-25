package twitch

import (
	"github.com/rs/zerolog/log"
	"github.com/streamcord/helix/v2"
	"time"
	"twproxy/dogstatsd"
)

// EnsureAppAccessToken requests an App Access Token and creates a background task to automatically refresh the token
// before it expires.
func (c *Client) EnsureAppAccessToken() {
	// Request and set an app access token
	creds := requestToken(c.Service, c.T)

	// Start the oauth token refresh service
	go func() {
		for {
			// Convert expiry time to time.Duration, and add a safety of 60 seconds to allow for token refresh
			expiresIn := time.Duration(creds.ExpiresIn-60) * time.Second
			log.Info().
				Str("service", c.Service).
				Msgf(
					"Access token will expire in %s (~%v days)",
					expiresIn.String(),
					int(expiresIn.Hours())/24)

			time.Sleep(expiresIn)

			creds = requestToken(c.Service, c.T)
		}
	}()
}

func requestToken(svcName string, c *helix.Client) helix.AccessCredentials {
	start := time.Now()
	res, err := c.RequestAppAccessToken([]string{})
	d := time.Now().Sub(start)
	go dogstatsd.LogTwitchRequest(dogstatsd.RouteGetOauthToken, svcName, res.ResponseCommon, err, d)

	if err != nil {
		log.Fatal().
			Str("service", svcName).
			Err(err).
			Msg("Failed to request app access token")
	} else if res.StatusCode > 399 {
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
