package config

import (
	"fmt"

	"github.com/b-open/jobbuzz/pkg/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Db gorm.DB

func GetDb(dbConfig DbConfig) (*gorm.DB, error) {
	// TODO: 12-Factor the connection string
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbConfig.Username,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.Database,
	)
	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}

func MigrateDb(db *gorm.DB) error {
	return db.AutoMigrate(&model.Job{})
}
