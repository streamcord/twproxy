package main

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"twproxy/config"
	"twproxy/dogstatsd"
	"twproxy/twitch"
)

func main() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if gin.IsDebugging() {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout}).With().Caller().Logger()
	}

	config.LoadConfig("config.yml")

	dogstatsd.New()
	twitch.CreateClientsFromServices(config.GlobalConfig.Services)

	r := NewRouter()
	r.Run("0.0.0.0:8181")
}
