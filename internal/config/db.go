package config

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GetDb() *gorm.DB {
	dsn := "host=localhost user=jobbuzz password=secret dbname=jobbuzz port=5532 sslmode=disable TimeZone=Asia/Brunei"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Fail to connect to DB", err)
	}

	return db
}
