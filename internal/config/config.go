package config

import "github.com/nix-united/golang-echo-boilerplate/internal/slogx"

type Config struct {
	Logger slogx.Config
	Auth   AuthConfig
	OAuth  OAuthConfig
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
	AccessSecret  string `env:"ACCESS_SECRET"`
	RefreshSecret string `env:"REFRESH_SECRET"`
}

type OAuthConfig struct {
	ClientID string `env:"OPEN_ID_CLIENT_ID"`
}

type HTTPConfig struct {
	Host       string `env:"HOST"`
	Port       string `env:"PORT"`
	ExposePort string `env:"EXPOSE_PORT"`
}
