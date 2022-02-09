package main

import (
	"log"

	"github.com/b-open/jobbuzz/internal/config"
)

func main() {

	configuration, err := config.LoadConfig("../../")

	if err != nil {
		log.Fatal("Fail to load db config", err)
	}

	dbConfig := configuration.DbConfig

	db, err := config.GetDb(dbConfig)

	if err != nil {
		log.Fatal("Fail to get db connection", err)
	}

	err = config.MigrateDb(db)

	if err != nil {
		log.Fatal("Fail to migrate db", err)
	}
}
