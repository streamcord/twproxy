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

	e := h.Group("/eventsub")
	e.DELETE("/subscriptions", eventsub.DeleteEventSubSubscription)
	e.GET("/subscriptions", eventsub.GetEventSubSubscriptions)
	e.POST("/subscriptions", eventsub.CreateEventSubSubscription)

	s := r.Group("/spyglass")
	s.POST("/callback", spyglass.PostCallback)

	return r
}
