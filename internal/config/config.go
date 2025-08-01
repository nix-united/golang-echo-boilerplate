package config

import (
	"time"
)

type Config struct {
	Logger LogConfig
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
	AccessTokenDuration  time.Duration `env:"ACCESS_SECRET_DURATION" envDefault:"2h"`
	RefreshTokenDuration time.Duration `env:"REFRESH_SECRET_DURATION" envDefault:"168h"`
	AccessSecret         string        `env:"ACCESS_SECRET"`
	RefreshSecret        string        `env:"REFRESH_SECRET"`
}

type OAuthConfig struct {
	ClientID string `env:"OPEN_ID_CLIENT_ID"`
}

type HTTPConfig struct {
	Host       string `env:"HOST"`
	Port       string `env:"PORT"`
	ExposePort string `env:"EXPOSE_PORT"`
}

type LogConfig struct {
	Application string `env:"LOG_APPLICATION"`

	// File represents path to file where store logs. Used [os.Stdout] if empty.
	File string `env:"LOG_FILE"`

	// One of: "DEBUG", "INFO", "WARN", "ERROR". Default: "DEBUG".
	Level string `env:"LOG_LEVEL" envDefault:"DEBUG"`

	// Add source code position to messages.
	AddSource bool `env:"LOG_ADD_SOURCE"`
}
