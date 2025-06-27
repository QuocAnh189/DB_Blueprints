package config

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"
)

var (
	cfg Config
)

const (
	ProductionEnv      = "production" //production or development
	DatabaseTimeout    = time.Second * 5
	ProductCachingTime = time.Minute * 1
)

type Config struct {
	HTTP_PORT   string `mapstructure:"HTTP_PORT"`
	DB_DRIVER   string `mapstructure:"DB_DRIVER"`
	DB_USER     string `mapstructure:"DB_USER"`
	DB_PASSWORD string `mapstructure:"DB_PASSWORD"`
	DB_HOST     string `mapstructure:"DB_HOST"`
	DB_PORT     string `mapstructure:"DB_PORT"`
	DB_NAME     string `mapstructure:"DB_NAME"`
}

func LoadConfig() *Config {
	viper.AutomaticEnv()

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		viper.SetConfigType("env")
		if err := viper.ReadInConfig(); err != nil {
			fmt.Printf("Error loading configuration file: %v", err)
		}
	}

	cfg = Config{
		HTTP_PORT:   viper.GetString("HTTP_PORT"),
		DB_DRIVER:   viper.GetString("DB_DRIVER"),
		DB_USER:     viper.GetString("DB_USER"),
		DB_PASSWORD: viper.GetString("DB_PASSWORD"),
		DB_HOST:     viper.GetString("DB_HOST"),
		DB_PORT:     viper.GetString("DB_PORT"),
		DB_NAME:     viper.GetString("DB_NAME"),
	}

	return &cfg
}

func GetConfig() *Config {
	return &cfg
}
