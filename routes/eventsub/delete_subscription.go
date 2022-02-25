package eventsub

import (
	"github.com/gin-gonic/gin"
	"github.com/nicklaw5/helix"
	"net/http"
	"time"
	"twproxy/dogstatsd"
)

// DeleteEventSubSubscription - Proxy of https://dev.twitch.tv/docs/api/reference#delete-eventsub-subscription
func DeleteEventSubSubscription(c *gin.Context) {
	t := c.MustGet("helix").(*helix.Client)

	start := time.Now()
	res, err := t.RemoveEventSubSubscription(c.Query("id"))
	d := time.Now().Sub(start)
	go dogstatsd.LogTwitchRequest(dogstatsd.RouteDeleteEventSubSubscription, c.GetHeader("Client-ID"), res.ResponseCommon, err, d)

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

	c.Status(http.StatusNoContent)
}
