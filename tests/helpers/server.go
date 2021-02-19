package helpers

import (
	"echo-demo-project/config"
	"echo-demo-project/server"

	"github.com/labstack/echo/v4"
)

func NewServer() *server.Server {
	_, redisClient := InitRedis()
	s := &server.Server{
		Echo: echo.New(),
		DB:   Init(),
		Redis: redisClient,
		Config: config.NewConfig(),
	}

	return s
}
