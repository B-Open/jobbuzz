package main

import (
	"os"

	"github.com/b-open/jobbuzz/internal/config"
	"github.com/b-open/jobbuzz/pkg/controller"
	"github.com/b-open/jobbuzz/pkg/middleware"
	"github.com/b-open/jobbuzz/pkg/service"
	"github.com/gin-gonic/gin"
	"github.com/mattn/go-isatty"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if gin.IsDebugging() {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	isTerm := isatty.IsTerminal(os.Stdout.Fd())
	if isTerm {
		log.Logger = log.Output(
			zerolog.ConsoleWriter{
				Out:     os.Stderr,
				NoColor: false,
			},
		)
	}

	configuration, err := config.LoadConfig("../../")
	if err != nil {
		log.Fatal().Err(err).Msg("Fail to load db config")
	}

	db, err := configuration.GetDb()
	if err != nil {
		log.Fatal().Err(err).Msg("Fail to get db connection")
	}

	service := service.Service{DB: db}
	controller := controller.Controller{Service: &service}

	r := gin.New()

	r.Use(middleware.SetLogger(isTerm))
	r.Use(gin.Recovery())

	r.GET("/api/jobs", controller.GetJobs)

	r.Run()
}
