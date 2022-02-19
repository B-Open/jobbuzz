package main

import (
	"os"

	"github.com/b-open/jobbuzz/internal/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(
		zerolog.ConsoleWriter{
			Out:     os.Stderr,
			NoColor: false,
		},
	)

	configuration, err := config.LoadConfig("../../")
	if err != nil {
		log.Fatal().Err(err).Msg("Fail to load db config")
	}

	db, err := configuration.GetDb()
	if err != nil {
		log.Fatal().Err(err).Msg("Fail to get db connection")
	}

	err = config.MigrateDb(db)
	if err != nil {
		log.Fatal().Err(err).Msg("Fail to migrate db")
	}

	log.Info().Msg("Migration completed")
}
