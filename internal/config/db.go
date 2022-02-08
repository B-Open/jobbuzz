package config

import (
	"log"

	"github.com/b-open/jobbuzz/pkg/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Db gorm.DB

// initialise DB once
func InitDb() {
	// TODO: 12-Factor the connection string
	dsn := "jobbuzz:secret@tcp(127.0.0.1:3306)/jobbuzz?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

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
