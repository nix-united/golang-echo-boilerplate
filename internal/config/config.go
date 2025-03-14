package config

type Config struct {
	Auth AuthConfig
	DB   DBConfig
	HTTP HTTPConfig
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

type HTTPConfig struct {
	Host       string `env:"HOST"`
	Port       string `env:"PORT"`
	ExposePort string `env:"EXPOSE_PORT"`
}
