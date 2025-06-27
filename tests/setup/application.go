package setup

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"os"
	"time"

	"github.com/testcontainers/testcontainers-go/modules/compose"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	applicationContainerName = "echo_app"
	applicationHost          = "localhost"
	applicationExposedPort   = "7788"
)

func SetupApplication(ctx context.Context) (_ *url.URL, _ func(context.Context) error, err error) {
	// If we attempt to run all tests at once, we may encounter an error
	// "Error response from daemon. No such container" due to the use of ryuk in Testcontainers.
	// Therefore, we need to disable it by setting an environment variable.
	if err := os.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true"); err != nil {
		return nil, nil, fmt.Errorf("set env to disable ryuk for testcontainers: %w", err)
	}

	dockerCompose, err := compose.NewDockerCompose("../../compose.yml")
	if err != nil {
		return nil, nil, fmt.Errorf("new docker compose: %w", err)
	}

	dockerCompose.
		WithEnv(map[string]string{
			"HOST":             applicationHost,
			"PORT":             applicationExposedPort,
			"DB_USER":          "local_user",
			"DB_PASSWORD":      "password",
			"DB_DRIVER":        "mysql",
			"DB_NAME":          "echo_example_docker",
			"DB_HOST":          "echo_mysql",
			"DB_PORT":          "3306",
			"COMPOSE_USER_ID":  "999",
			"COMPOSE_GROUP_ID": "999",
			"EXPOSE_PORT":      applicationExposedPort,
			"EXPOSE_DB_PORT":   "33060",
			"ACCESS_SECRET":    "access_secret",
			"REFRESH_SECRET":   "refresh_secret",
		}).
		WaitForService(
			applicationContainerName,
			wait.ForAll(wait.ForHTTP("/health")).WithDeadline(5*time.Minute),
		)

	if err := dockerCompose.Up(ctx); err != nil {
		return nil, nil, fmt.Errorf("docker compose up: %w", err)
	}

	shutdown := func(ctx context.Context) error {
		application, err := dockerCompose.ServiceContainer(ctx, applicationContainerName)
		if err != nil {
			return fmt.Errorf("get application container: %w", err)
		}

		l, err := application.Logs(ctx)
		if err != nil {
			return fmt.Errorf("get logs: %w", err)
		}

		r, err := io.ReadAll(l)
		if err != nil {
			return fmt.Errorf("read logs: %w", err)
		}

		fmt.Println(string(r))

		if err := dockerCompose.Down(ctx); err != nil {
			return fmt.Errorf("down docker compose: %w", err)
		}

		return nil
	}

	appliactionURL, err := url.Parse(fmt.Sprintf("http://%s:%s", applicationHost, applicationExposedPort))
	if err != nil {
		return nil, nil, fmt.Errorf("parse application url: %w", err)
	}

	return appliactionURL, shutdown, nil
}
