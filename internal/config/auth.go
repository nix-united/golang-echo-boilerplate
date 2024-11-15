package config

import "os"

type AuthConfig struct {
	AccessSecret  string
	RefreshSecret string
}

func LoadAuthConfig() AuthConfig {
	return AuthConfig{
		AccessSecret:  os.Getenv("ACCESS_SECRET"),
		RefreshSecret: os.Getenv("REFRESH_SECRET"),
	}
}
