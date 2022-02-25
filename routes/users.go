package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nicklaw5/helix"
	"time"
	"twproxy/dogstatsd"
)

// GetUsers - Proxy of https://dev.twitch.tv/docs/api/reference#get-users
func GetUsers(c *gin.Context) {
	t := c.MustGet("helix").(*helix.Client)

	start := time.Now()
	res, err := t.GetUsers(&helix.UsersParams{
		IDs:    c.QueryArray("id"),
		Logins: c.QueryArray("login"),
	})
	d := time.Now().Sub(start)
	go dogstatsd.LogTwitchRequest(dogstatsd.RouteGetUsers, c.GetHeader("Client-ID"), res.ResponseCommon, err, d)

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
