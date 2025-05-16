package setup

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	mysqlImage    = "mysql:9.3.0"
	mysqlDatabase = "db_name"
	mysqlUsername = "username"
	mysqlPassword = "password"
	mysqlPort     = "3306"
	mysqlHost     = "localhost"
)

type MySQLConfig struct {
	User        string
	Password    string
	Host        string
	ExposedPort string
	LocalPort   string
	Name        string
}

func SetupMySQL(ctx context.Context) (_ MySQLConfig, _ func(ctx context.Context) error, err error) {
	container, err := mysql.Run(
		ctx,
		mysqlImage,
		mysql.WithDatabase(mysqlDatabase),
		mysql.WithUsername(mysqlUsername),
		mysql.WithPassword(mysqlPassword),
		testcontainers.WithWaitStrategyAndDeadline(
			time.Minute,
			wait.ForLog(fmt.Sprintf(
				"/usr/sbin/mysqld: ready for connections. Version: '9.3.0'  socket: '/var/run/mysqld/mysqld.sock'  port: %s  MySQL Community Server - GPL.",
				mysqlPort,
			)),
		),
	)
	if err != nil {
		return MySQLConfig{}, nil, fmt.Errorf("run mysql container: %w", err)
	}

	shutdown := func(ctx context.Context) error {
		if err := container.Terminate(ctx); err != nil {
			return fmt.Errorf("terminate mysql container: %w", err)
		}

		return nil
	}

	defer func() {
		if err == nil {
			return
		}

		if errShutdown := shutdown(ctx); errShutdown != nil {
			err = errors.Join(err, errShutdown)
		}
	}()

	port, err := container.MappedPort(ctx, mysqlPort+"/tcp")
	if err != nil {
		return MySQLConfig{}, nil, fmt.Errorf("get mysql exposed port: %w", err)
	}

	config := MySQLConfig{
		User:        mysqlUsername,
		Password:    mysqlPassword,
		Host:        mysqlHost,
		ExposedPort: port.Port(),
		LocalPort:   mysqlPort,
		Name:        mysqlDatabase,
	}

	return config, shutdown, nil
}
