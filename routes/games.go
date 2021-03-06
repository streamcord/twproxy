package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/streamcord/helix/v2"
	"time"
	"twproxy/dogstatsd"
	"twproxy/twitch"
)

// GetGames - Proxy of https://dev.twitch.tv/docs/api/reference#get-games
func GetGames(c *gin.Context) {
	t := c.MustGet("helix").(*twitch.Client).T

	start := time.Now()
	res, err := t.GetGames(&helix.GamesParams{
		IDs:   c.QueryArray("id"),
		Names: c.QueryArray("name"),
	})
	d := time.Now().Sub(start)
	go dogstatsd.LogTwitchRequest(dogstatsd.RouteGetGames, c.GetHeader("Client-ID"), res.ResponseCommon, err, d)

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
