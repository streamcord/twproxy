package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nicklaw5/helix"
)

// GetGames - Proxy of https://dev.twitch.tv/docs/api/reference#get-games
func GetGames(c *gin.Context) {
	t := c.MustGet("helix").(*helix.Client)

	res, err := t.GetGames(&helix.GamesParams{
		IDs:   c.QueryArray("id"),
		Names: c.QueryArray("name"),
	})
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
