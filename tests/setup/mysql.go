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
	mysqlImage         = "mysql:9.3.0"
	mysqlDatabase      = "db_name"
	mysqlUsername      = "username"
	mysqlPassword      = "password"
	mysqlPort          = "3306"
	mysqlHost          = "localhost"
	mysqlContainerName = "golang_echo_boilerplate_mysql_db"
)

type MySQLConfig struct {
	User          string
	Password      string
	Host          string
	ExposedPort   string
	LocalPort     string
	Name          string
	ContainerName string
}

func SetupMySQL(ctx context.Context, networks []string) (MySQLConfig, func(ctx context.Context) error, error) {
	containerLogsConsumer := newContainerLogsConsumer(mysqlContainerName)

	container, err := mysql.Run(
		ctx,
		mysqlImage,
		mysql.WithDatabase(mysqlDatabase),
		mysql.WithUsername(mysqlUsername),
		mysql.WithPassword(mysqlPassword),
		testcontainers.WithWaitStrategyAndDeadline(
			time.Minute,
			wait.ForLog("X Plugin ready for connections. Bind-address: '::' port: 33060, socket: /var/run/mysqld/mysqlx.sock"),
		),
		testcontainers.CustomizeRequest(testcontainers.GenericContainerRequest{
			ContainerRequest: testcontainers.ContainerRequest{
				Name:     mysqlContainerName,
				Networks: networks,
				LogConsumerCfg: &testcontainers.LogConsumerConfig{
					Consumers: []testcontainers.LogConsumer{containerLogsConsumer},
				},
			},
		}),
	)
	if err != nil {
		// Print logs for container bootstrap failures, such as condition wait timeouts.
		containerLogsConsumer.Print()
		return MySQLConfig{}, nil, fmt.Errorf("run mysql container: %w", err)
	}

	shutdown := func(ctx context.Context) error {
		if err := container.Terminate(ctx); err != nil {
			return fmt.Errorf("terminate mysql container: %w", err)
		}

		// Print container logs after tests complete during shutdown.
		containerLogsConsumer.Print()

		return nil
	}

	// This defer function exists to shutdown the container if any error occurs in next steps.
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
		User:          mysqlUsername,
		Password:      mysqlPassword,
		Host:          mysqlHost,
		ExposedPort:   port.Port(),
		LocalPort:     mysqlPort,
		Name:          mysqlDatabase,
		ContainerName: mysqlContainerName,
	}

	return config, shutdown, nil
}
