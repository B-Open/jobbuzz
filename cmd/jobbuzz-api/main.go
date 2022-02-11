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

	r.Use(logger.SetLogger(
		logger.WithLogger(func(c *gin.Context, out io.Writer, latency time.Duration) zerolog.Logger {
			logger := zerolog.New(out)
			if isTerm {
				logger = logger.Output(
					zerolog.ConsoleWriter{
						Out:     out,
						NoColor: false,
					},
				)
			}
			logger = logger.With().
				Timestamp().
				Int("status", c.Writer.Status()).
				Str("method", c.Request.Method).
				Str("path", c.Request.URL.Path).
				Str("ip", c.ClientIP()).
				Dur("latency", latency).
				Str("user_agent", c.Request.UserAgent()).
				Logger()
			return logger
		}),
	))
	r.Use(gin.Recovery())

	r.GET("/api/jobs", controller.GetJobs)

	r.Run()
}
