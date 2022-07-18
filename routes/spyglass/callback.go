package spyglass

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/streamcord/helix/v2"
	"io/ioutil"
)

// EventSubWebhookCallbackVerificationEvent ...
type EventSubWebhookCallbackVerificationEvent struct {
	Challenge    string                     `json:"challenge"`
	Subscription helix.EventSubSubscription `json:"subscription"`
}

// EventSubRevocationEvent ...
type EventSubRevocationEvent struct {
	Subscription helix.EventSubSubscription `json:"subscription"`
}

// EventSubNotificationEvent ...
type EventSubNotificationEvent struct {
	Subscription helix.EventSubSubscription `json:"subscription"`
	Event        struct {
		BroadcasterUserID    string     `json:"broadcaster_user_id"`
		BroadcasterUserLogin string     `json:"broadcaster_user_login"`
		StartedAt            helix.Time `json:"started_at,omitempty"`
		Type                 string     `json:"type,omitempty"`
	} `json:"event"`
}

// PostCallback - Handles incoming EventSub webhooks
func PostCallback(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)

	secret := "Password5555"

	if !helix.VerifyEventSubNotification(secret, c.Request.Header, string(body)) {
		c.AbortWithStatusJSON(401, gin.H{
			"error": "Invalid signature",
		})
		return
	}

	msgType := c.GetHeader("Twitch-Eventsub-Message-Type")
	switch msgType {
	case "webhook_callback_verification":
		// Process a verification request from Twitch
		// https://dev.twitch.tv/docs/eventsub/handling-webhook-events#responding-to-a-challenge-request
		var res EventSubWebhookCallbackVerificationEvent
		err := json.Unmarshal(body, &res)
		if err != nil {
			log.Warn().Err(err).Msgf("Failed to verify subscription %s", res.Subscription.ID)
			c.AbortWithStatusJSON(400, gin.H{
				"error": "Unable to parse request body: " + err.Error(),
			})
			return
		}

		log.Info().Msgf("Verified subscription %s", res.Subscription.ID)
		c.String(200, res.Challenge)
		break
	case "revocation":
		// Notification of subscription revocation
		// https://dev.twitch.tv/docs/eventsub/handling-webhook-events#revoking-your-subscription
		// TODO: Check reason for revocation and recreate sub or delete notifications
		var res EventSubRevocationEvent
		err := json.Unmarshal(body, &res)
		if err != nil {
			log.Warn().Err(err).Msgf("Failed to revoke subscription %s", res.Subscription.ID)
			c.AbortWithStatusJSON(400, gin.H{
				"error": "Unable to parse request body: " + err.Error(),
			})
			return
		}

		log.Info().Str("reason", res.Subscription.Status).Msgf("Revoked subscription %s", res.Subscription.ID)
		c.String(201, "")
		break
	case "notification":
		// Normal event received
		// https://dev.twitch.tv/docs/eventsub/handling-webhook-events#processing-an-event
		var res EventSubNotificationEvent
		err := json.Unmarshal(body, &res)
		if err != nil {
			log.Warn().Err(err).Msgf("Failed to process subscription notification for %s", res.Subscription.ID)
			c.AbortWithStatusJSON(400, gin.H{
				"error": "Unable to parse request body: " + err.Error(),
			})
			return
		}
		log.Info().Interface("data", res).Msgf("Received event %s for subscription %s", res.Subscription.Type, res.Subscription.ID)
		c.String(201, "")
		break
	default:
		log.Warn().
			Str("type", msgType).
			Msg("Received unknown message type")
		c.String(501, "")
	}
}
