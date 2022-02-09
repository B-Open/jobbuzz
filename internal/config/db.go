package config

import (
	"fmt"

	"github.com/b-open/jobbuzz/pkg/model"
	gomysql "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Db gorm.DB

func (configuration *Configuration) GetDb() (*gorm.DB, error) {
	dsn := configuration.formatDsn().FormatDSN()
	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}

func MigrateDb(db *gorm.DB) error {
	return db.AutoMigrate(&model.Job{})
}

func (configuration *Configuration) formatDsn() *gomysql.Config {
	dbConfig := configuration.DbConfig

	mysqlConfig := gomysql.NewConfig()
	mysqlConfig.Net = "tcp"
	mysqlConfig.Addr = fmt.Sprintf("%s:%s", dbConfig.Host, dbConfig.Port)
	mysqlConfig.User = dbConfig.Username
	mysqlConfig.Passwd = dbConfig.Password
	mysqlConfig.DBName = dbConfig.Database
	mysqlConfig.ParseTime = true

	return mysqlConfig
}
