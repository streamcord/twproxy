package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/streamcord/helix/v2"
	"time"
	"twproxy/dogstatsd"
	"twproxy/twitch"
	"twproxy/utils"
)

// GetUsersFollows - Proxy of https://dev.twitch.tv/docs/api/reference#get-users-follows
func GetUsersFollows(c *gin.Context) {
	t := c.MustGet("helix").(*twitch.Client).T

	start := time.Now()
	res, err := t.GetUsersFollows(&helix.UsersFollowsParams{
		After:  c.Query("after"),
		First:  utils.QueryInt(c, "first", 20),
		FromID: c.Query("from_id"),
		ToID:   c.Query("to_id"),
	})
	d := time.Now().Sub(start)
	go dogstatsd.LogTwitchRequest(dogstatsd.RouteGetUsersFollows, c.GetHeader("Client-ID"), res.ResponseCommon, err, d)

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
