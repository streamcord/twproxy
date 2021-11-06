package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nicklaw5/helix"
	"twproxy/utils"
)

// GetClips - Proxy of https://dev.twitch.tv/docs/api/reference#get-clips
func GetClips(c *gin.Context) {
	t := c.MustGet("helix").(*helix.Client)

	res, err := t.GetClips(&helix.ClipsParams{
		BroadcasterID: c.Query("broadcaster_id"),
		GameID:        c.Query("game_id"),
		IDs:           c.QueryArray("id"),
		First:         utils.QueryInt(c, "first", 20),
		After:         c.Query("after"),
		Before:        c.Query("before"),
		EndedAt:       utils.QueryTime(c, "ended_at"),
		StartedAt:     utils.QueryTime(c, "started_at"),
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