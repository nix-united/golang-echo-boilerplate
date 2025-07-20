package config

import (
	"time"

	"github.com/nix-united/golang-echo-boilerplate/internal/slogx"
)

type Config struct {
	Logger slogx.Config
	Auth   AuthConfig
	DB     DBConfig
	HTTP   HTTPConfig
}

type DBConfig struct {
	User     string `env:"DB_USER"`
	Password string `env:"DB_PASSWORD"`
	Driver   string `env:"DB_DRIVER"`
	Name     string `env:"DB_NAME"`
	Host     string `env:"DB_HOST"`
	Port     string `env:"DB_PORT"`
}

type AuthConfig struct {
	AccessTokenDuration  time.Duration `env:"ACCESS_SECRET_DURATION" envDefault:"2h"`
	RefreshTokenDuration time.Duration `env:"REFRESH_SECRET_DURATION" envDefault:"168h"`
	AccessSecret         string        `env:"ACCESS_SECRET"`
	RefreshSecret        string        `env:"REFRESH_SECRET"`
}

type HTTPConfig struct {
	Host       string `env:"HOST"`
	Port       string `env:"PORT"`
	ExposePort string `env:"EXPOSE_PORT"`
}
