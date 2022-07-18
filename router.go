package main

import (
	"github.com/gin-gonic/gin"
	"twproxy/middleware"
	"twproxy/routes"
	"twproxy/routes/eventsub"
	"twproxy/routes/spyglass"
)

// NewRouter ...
func NewRouter() *gin.Engine {
	r := gin.New()
	r.Use(middleware.LoggerMiddleware) // Register last

	h := r.Group("/helix")
	h.Use(middleware.AuthMiddleware)

	h.GET("/clips", routes.GetClips)
	h.GET("/games", routes.GetGames)
	h.GET("/streams", routes.GetStreams)
	h.GET("/users", routes.GetUsers)
	h.GET("/users/follows", routes.GetUsersFollows)

	// /helix/eventsub/...
	e := h.Group("/eventsub")
	e.DELETE("/subscriptions", eventsub.DeleteEventSubSubscription)
	e.GET("/subscriptions", eventsub.GetEventSubSubscriptions)
	e.POST("/subscriptions", eventsub.CreateEventSubSubscription)

	// /spyglass/...
	s := r.Group("/spyglass")
	s.POST("/callback", spyglass.PostCallback)

	// /spyglass/guilds/:guild_id/channels/:channel_id/...
	g := s.Group("/guilds/:guild_id/channels/:channel_id")
	g.POST("/notifications", spyglass.CreateNotification)
	g.DELETE("/notifications/:notification_id", spyglass.DeleteNotification)
	g.PATCH("/notifications/:notification_id", spyglass.UpdateNotification)

	return r
}
