package config

import "github.com/nix-united/golang-echo-boilerplate/internal/db"

type Config struct {
	Auth AuthConfig
	DB   db.Config
	HTTP HTTPConfig
}

type AuthConfig struct {
	AccessSecret  string `env:"ACCESS_SECRET"`
	RefreshSecret string `env:"REFRESH_SECRET"`
}

type HTTPConfig struct {
	Host       string `env:"HOST"`
	Port       string `env:"PORT"`
	ExposePort string `env:"EXPOSE_PORT"`
}
