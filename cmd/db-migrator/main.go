package main

import "github.com/b-open/jobbuzz/internal/config"

func main() {

	config.InitDb()

	db := config.GetDb()

	config.MigrateDb(db)
}
