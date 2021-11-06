package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/nicklaw5/helix"
	"github.com/rs/zerolog/log"
	"strconv"
)

// QueryInt returns the keyed url query value if it exists, otherwise it returns the value of `def`.
func QueryInt(c *gin.Context, key string, def int) int {
	raw, ok := c.GetQuery(key)
	if ok {
		var err error
		def, err = strconv.Atoi(raw)
		if err != nil {
			def = 20
		}
	}
	return def
}

// QueryTime returns the keyed url query value as helix.Time if it exists.
func QueryTime(c *gin.Context, key string) helix.Time {
	raw, ok := c.GetQuery(key)
	if ok {
		var t helix.Time
		err := (*helix.Time).UnmarshalJSON(&t, []byte(raw))
		if err != nil {
			log.Warn().Err(err).Msg("Failed to unmarshal to helix.Time")
			return helix.Time{}
		}
		return t
	}
	return helix.Time{}
}
