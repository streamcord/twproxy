package spyglass

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
	"twproxy/utils"
)

// CreateNotification - The client should send the notification data as a JSON-encoded payload.
//
// Since the request will only come from an authoritative application (i.e. the dashboard), the option of whether
// to enable pro features should be sent in the `guild_pro_level` query parameter.
//
// 0 = free bot, no pro features. 1 = pro enabled.
//
// This option should only be used to calculate the number of notifications allowed for the guild, and not whether
// certain premium-only notification features should be enabled because spyglass-receiver determines this.
//
// This endpoint also assumes that the guild and channel IDs are valid.
func CreateNotification(c *gin.Context) {
	var body utils.SpyglassNotificationPayload
	if err := c.ShouldBindJSON(&body); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	notif := utils.SpyglassNotification{
		ID:                          utils.GenerateID("12345678"),
		SpyglassNotificationPayload: body,
		Channel:                     c.Param("channel_id"),
		CreatedAt:                   primitive.DateTime(time.Now().Unix()),
		Guild:                       c.Param("guild_id"),
		StreamerID:                  "",
		UpdatedAt:                   primitive.DateTime(time.Now().Unix()),
	}

	c.JSON(200, notif)
}
