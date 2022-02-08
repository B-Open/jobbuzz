package config

import (
	"log"

	"github.com/b-open/jobbuzz/pkg/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Db gorm.DB

// initialise DB once
func InitDb() {
	// TODO: 12-Factor the connection string
	dsn := "host=localhost user=jobbuzz password=secret dbname=jobbuzz port=5432 sslmode=disable TimeZone=Asia/Brunei"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Fail to connect to DB", err)
	}

	Db = *db
}

// reusing the same DB instance
func GetDb() *gorm.DB {
	return &Db
}

func MigrateDb(db *gorm.DB) {
	err := db.AutoMigrate(&model.Job{})

	if err != nil {
		log.Fatal(err)
	}
}
