package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	DBDriver   string `mapstructure:"db_driver"`
	DBSource   string `mapstructure:"db_source"`
	ServerPort string `mapstructure:"server_port"`
	LogLevel   string `mapstructure:"log_level"`
}

var AppConfig Config

func LoadConfig(configPath string) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configPath)
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("No config file found: %v, using env vars only", err)
	}

	if err := viper.Unmarshal(&AppConfig); err != nil {
		log.Fatalf("Failed to unmarshal config: %v", err)
	}

	if AppConfig.ServerPort == "" {
		AppConfig.ServerPort = ":8080"
	}
	if AppConfig.DBDriver == "" {
		AppConfig.DBDriver = "sqlite3"
	}
	if AppConfig.DBSource == "" {
		AppConfig.DBSource = "trader.db"
	}
	if AppConfig.LogLevel == "" {
		AppConfig.LogLevel = "info"
	}
}
