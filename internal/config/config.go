package config

import (
	"github.com/spf13/viper"
)

type Configuration struct {
	DbConfig DbConfig
}

type DbConfig struct {
	Host     string `mapstructure:"DB_HOST"`
	Port     string `mapstructure:"DB_PORT"`
	Username string `mapstructure:"DB_USERNAME"`
	Password string `mapstructure:"DB_PASSWORD"`
	Database string `mapstructure:"DB_DATABASE"`
}

var configuration Configuration

func LoadConfig(path string) (*Configuration, error) {

	viper.AddConfigPath(path)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	// load the db configuration
	dbConfig, err := loadDbConfig(path)

	if err != nil {
		return nil, err
	}

	configuration.DbConfig = *dbConfig

	return &configuration, nil
}

func loadDbConfig(path string) (*DbConfig, error) {

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	var dbConfig DbConfig

	err = viper.Unmarshal(&dbConfig)
	if err != nil {
		return nil, err
	}

	return &dbConfig, nil
}
