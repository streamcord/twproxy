package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nicklaw5/helix"
	"time"
	"twproxy/dogstatsd"
	"twproxy/utils"
)

// GetStreams - Proxy of https://dev.twitch.tv/docs/api/reference#get-streams
func GetStreams(c *gin.Context) {
	t := c.MustGet("helix").(*helix.Client)

	start := time.Now()
	res, err := t.GetStreams(&helix.StreamsParams{
		After:      c.Query("after"),
		Before:     c.Query("before"),
		First:      utils.QueryInt(c, "first", 20),
		GameIDs:    c.QueryArray("game_id"),
		Language:   c.QueryArray("language"),
		Type:       c.Query("type"),
		UserIDs:    c.QueryArray("user_id"),
		UserLogins: c.QueryArray("user_login"),
	})
	d := time.Now().Sub(start)
	go dogstatsd.LogTwitchRequest(dogstatsd.RouteGetStreams, c.GetHeader("Client-ID"), res.ResponseCommon, err, d)

	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{
			"error":   "Internal Server Error",
			"message": err.Error(),
			"status":  500,
		})
		return
	} else if res.StatusCode != 200 {
		c.AbortWithStatusJSON(res.StatusCode, gin.H{
			"error":   res.Error,
			"message": res.ErrorMessage,
			"status":  res.ErrorStatus,
		})
		return
	}

	c.JSON(200, res.Data)
}
