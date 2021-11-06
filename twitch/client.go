package twitch

import (
	"github.com/nicklaw5/helix"
	"github.com/rs/zerolog/log"
	"net/http"
	"time"
	"twproxy/config"
)

// GlobalClients - map of (service name: client object)
var GlobalClients map[string]*helix.Client

// NewClientFromService ...
func NewClientFromService(name string, s config.Service) *helix.Client {
	// Create a custom HTTP client
	// Borrowed from https://www.loginradius.com/blog/async/tune-the-go-http-client-for-high-performance/
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = 100
	t.MaxConnsPerHost = 100
	t.MaxIdleConnsPerHost = 100

	httpClient := &http.Client{
		Timeout:   10 * time.Second,
		Transport: t,
	}

	// Create the helix client
	c, err := helix.NewClient(&helix.Options{
		ClientID:      s.ClientID,
		ClientSecret:  s.ClientSecret,
		HTTPClient:    httpClient,
		RateLimitFunc: ratelimitCallback,
		UserAgent:     "Go-http-client/1.1 (github.com/streamcord/twproxy v1.0, github.com/nicklaw5/helix v2.2.0)",
	})
	if err != nil {
		log.Fatal().
			Str("service", name).
			Err(err).
			Msg("Failed to create helix client")
	}

	// Make sure that there is always a valid app access token
	EnsureAppAccessToken(name, c)

	return c
}

// CreateClientsFromServices - Update GlobalClients from a list of config.Service.
func CreateClientsFromServices(s map[string]config.Service) {
	GlobalClients = make(map[string]*helix.Client, len(s))

	for name, svc := range s {
		log.Debug().Msgf("Creating helix client for service '%s'", name)
		GlobalClients[name] = NewClientFromService(name, svc)
	}
}