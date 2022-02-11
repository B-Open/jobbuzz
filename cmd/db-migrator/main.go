package main

import (
	"github.com/b-open/jobbuzz/internal/config"
	"github.com/rs/zerolog/log"
)

func main() {
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
