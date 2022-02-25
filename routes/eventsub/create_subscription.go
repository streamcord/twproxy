package eventsub

import (
	"github.com/gin-gonic/gin"
	"github.com/nicklaw5/helix"
	"time"
	"twproxy/dogstatsd"
)

// CreateEventSubSubscription - Proxy of https://dev.twitch.tv/docs/api/reference#create-eventsub-subscription
func CreateEventSubSubscription(c *gin.Context) {
	t := c.MustGet("helix").(*helix.Client)

	var sub helix.EventSubSubscription
	err := c.BindJSON(&sub)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
			"status":  400,
		})
		return
	}

	start := time.Now()
	res, err := t.CreateEventSubSubscription(&sub)
	d := time.Now().Sub(start)
	go dogstatsd.LogTwitchRequest(dogstatsd.RouteCreateEventSubSubscription, c.GetHeader("Client-ID"), res.ResponseCommon, err, d)
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
