package twitch

import (
	"github.com/streamcord/helix/v2"
)

// Client ...
type Client struct {
	Service string
	T       *helix.Client
}
