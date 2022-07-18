package utils

import (
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strconv"
)

// SpyglassNotificationID - Extension of primitive.Binary with JSON marshalling support
type SpyglassNotificationID primitive.Binary

// String ...
func (id SpyglassNotificationID) String() string {
	return base64.URLEncoding.EncodeToString(id.Data)
}

// MarshalJSON ...
func (id SpyglassNotificationID) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.String())
}

// SpyglassNotificationLastEventError ...
// https://github.com/streamcord/spyglass-receiver/blob/master/types/notification/notification.go#L12
type SpyglassNotificationLastEventError struct {
	Code    int    `bson:"code"`
	Message string `bson:"message"`
	Status  int    `bson:"status"`
}

// SpyglassNotificationLastEvent ...
// https://github.com/streamcord/spyglass-receiver/blob/master/types/notification/notification.go#L19
type SpyglassNotificationLastEvent struct {
	At                    primitive.DateTime                  `json:"at" bson:"at"`
	Error                 *SpyglassNotificationLastEventError `json:"error,omitempty" bson:"error,omitempty"`
	MessageID             string                              `json:"message_id" bson:"message_id"`
	NotifiedStreamOffline bool                                `json:"notified_stream_offline" bson:"notified_stream_offline"`
	RegexMatched          bool                                `json:"regex_matched" bson:"regex_matched"`
	StreamID              string                              `json:"stream_id" bson:"stream_id"`
}

// SpyglassNotificationRegexPattern ...
type SpyglassNotificationRegexPattern struct {
	Field     string `json:"field" binding:"regexfield"`
	MatchType bool   `json:"match_type" binding:"required"`
	Pattern   string `json:"pattern" binding:"required,lte=2000"`
}

// SpyglassNotificationRegex ...
// https://github.com/streamcord/spyglass-receiver/blob/master/types/notification/notification.go#L51
type SpyglassNotificationRegex struct {
	Logic    uint                               `json:"logic" binding:"lte=1"`
	Patterns []SpyglassNotificationRegexPattern `json:"patterns" binding:"regexpattern"`
}

// SpyglassNotificationPayload - Object to bind with request bodies for POST and PATCH requests to the Spyglass
// notifications endpoint.
//
// EmbedFlags: https://github.com/streamcord/spyglass-receiver/blob/master/constants/flags.go
// Valid range: 0-127
//
// StreamEndAction: https://github.com/streamcord/spyglass-receiver/blob/master/constants/stream_end_actions.go
// Valid range: 0-3
type SpyglassNotificationPayload struct {
	Cooldown         uint                      `json:"cooldown,omitempty" binding:"lte=10080" bson:"cooldown,omitempty"`
	CreatedBy        string                    `json:"created_by,omitempty" bson:"created_by,omitempty"`
	EmbedColor       string                    `json:"embed_color,omitempty" binding:"hex,omitempty" bson:"embed_color,omitempty"`
	EmbedFlags       uint                      `json:"embed_flags,omitempty" binding:"lt=128" bson:"embed_flags,omitempty"`
	ImageURL         string                    `json:"image_url,omitempty" binding:"optionalurl" bson:"image_url,omitempty"`
	Message          string                    `json:"message" binding:"required,lte=2000" bson:"message"`
	Regex            SpyglassNotificationRegex `json:"regex" bson:"regex"`
	StreamerName     string                    `json:"streamer_name" binding:"required,twitchusername" bson:"streamer_name"`
	StreamEndAction  uint                      `json:"stream_end_action,omitempty" binding:"lte=3" bson:"stream_end_action,omitempty"`
	StreamEndMessage string                    `json:"stream_end_message,omitempty" binding:"lte=2000" bson:"stream_end_message,omitempty"`
}

// SpyglassNotification complete object.
// https://github.com/streamcord/spyglass-receiver/blob/master/types/notification/notification.go#L59
type SpyglassNotification struct {
	SpyglassNotificationPayload
	ID         SpyglassNotificationID         `json:"_id" bson:"_id"`
	Channel    string                         `json:"channel" bson:"channel"`
	CreatedAt  primitive.DateTime             `json:"created_at" bson:"created_at"`
	Guild      string                         `json:"guild" bson:"guild"`
	LastEvent  *SpyglassNotificationLastEvent `json:"last_event,omitempty" bson:"last_event,omitempty"`
	StreamerID string                         `json:"streamer_id" bson:"streamer_id"`
	UpdatedAt  primitive.DateTime             `json:"updated_at" bson:"updated_at"`
}

// IDString converts a notification ID to a base64 string.
func (n SpyglassNotification) IDString() string {
	return base64.URLEncoding.EncodeToString(n.ID.Data)
}

// GenerateID - Generate a new Spyglass ID.
func GenerateID(streamerID string) SpyglassNotificationID {
	oldBytesFixed := [12]byte(primitive.NewObjectID())
	oldBytes := oldBytesFixed[:]

	idInt, _ := strconv.Atoi(streamerID)
	idBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(idBytes, uint64(idInt))

	newBytes := append(oldBytes, idBytes...)
	return SpyglassNotificationID{
		Data: newBytes,
	}
}
