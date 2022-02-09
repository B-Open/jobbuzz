package main

import (
	"log"

	"github.com/b-open/jobbuzz/internal/config"
)

func main() {

	db, err := config.GetDb()

	if err != nil {
		log.Fatal("Fail to get db connection", err)
	}

	err = config.MigrateDb(db)

	if err != nil {
		log.Fatal("Fail to migrate db", err)
	}
}
