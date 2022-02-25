package twitch

import (
	"github.com/nicklaw5/helix"
)

// Client ...
type Client struct {
	Service string
	T       *helix.Client
}
