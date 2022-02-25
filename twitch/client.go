package twitch

import (
	"github.com/nicklaw5/helix"
)

type Client struct {
	Service string
	T       *helix.Client
}
