package eventsub

import (
	"github.com/gin-gonic/gin"
	"github.com/streamcord/helix/v2"
	"time"
	"twproxy/dogstatsd"
	"twproxy/twitch"
)

// GetEventSubSubscriptions - Proxy of https://dev.twitch.tv/docs/api/reference#get-eventsub-subscriptions
func GetEventSubSubscriptions(c *gin.Context) {
	t := c.MustGet("helix").(*twitch.Client).T

	start := time.Now()
	res, err := t.GetEventSubSubscriptions(&helix.EventSubSubscriptionsParams{
		Status: c.Query("status"),
		Type:   c.Query("type"),
		After:  c.Query("after"),
	})
	d := time.Now().Sub(start)
	go func() {
		svc := c.GetHeader("Client-ID")
		dogstatsd.LogTwitchRequest(dogstatsd.RouteGetEventSubSubscriptions, svc, res.ResponseCommon, err, d)
		dogstatsd.LogEventSubSubscriptionCount(res, svc)
	}()

	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{
			"error":   "Internal Server Error",
			"message": err.Error(),
			"status":  500,
		})
		return
	} else if res.StatusCode > 399 {
		c.AbortWithStatusJSON(res.StatusCode, gin.H{
			"error":   res.Error,
			"message": res.ErrorMessage,
			"status":  res.ErrorStatus,
		})
		return
	}

	c.JSON(res.StatusCode, res.Data)
}
