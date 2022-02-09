package main

import (
	"log"

	"github.com/b-open/jobbuzz/internal/config"
)

func main() {

	dbCofig, err := config.LoadDbConfig("../../")

	if err != nil {
		log.Fatal("Fail to load db config", err)
	}

	db, err := config.GetDb(*dbCofig)

	if err != nil {
		log.Fatal("Fail to get db connection", err)
	}

	err = config.MigrateDb(db)

	if err != nil {
		log.Fatal("Fail to migrate db", err)
	}
}
