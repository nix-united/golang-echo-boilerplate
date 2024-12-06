package config

import (
	"log"

	"github.com/joho/godotenv"
)

type Config struct {
	Auth AuthConfig
	DB   DBConfig
	HTTP HTTPConfig
}

func NewConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	return &Config{
		Auth: LoadAuthConfig(),
		DB:   LoadDBConfig(),
		HTTP: LoadHTTPConfig(),
	}
}
