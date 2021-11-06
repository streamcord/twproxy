package main

import (
	"github.com/gin-gonic/gin"
	"twproxy/middleware"
	"twproxy/routes"
)

// NewRouter ...
func NewRouter() *gin.Engine {
	r := gin.Default()

	h := r.Group("/helix")
	h.Use(middleware.AuthMiddleware)

	h.GET("/clips", routes.GetClips)
	h.GET("/games", routes.GetGames)
	h.GET("/streams", routes.GetStreams)
	h.GET("/users", routes.GetUsers)
	h.GET("/users/follows", routes.GetUsersFollows)

	return r
}
