package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	DbConfig DbConfig
}

type DbConfig struct {
	Host     string `mapstructure:"DB_HOST"`
	Port     string `mapstructure:"DB_PORT"`
	Username string `mapstructure:"DB_USERNAME"`
	Password string `mapstructure:"DB_PASSWORD"`
	Database string `mapstructure:"DB_DATABASE"`
}

var config Config

func LoadConfig(path string) (*Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	var dbConfig DbConfig

	err = viper.Unmarshal(&dbConfig)
	if err != nil {
		return nil, err
	}

	config.DbConfig = dbConfig

	return &config, nil
}

func (c *Config) GetDbConfig() *DbConfig {
	return &c.DbConfig
}
