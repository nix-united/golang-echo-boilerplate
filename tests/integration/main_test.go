package integration

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"slices"
	"testing"

	"github.com/nix-united/golang-echo-boilerplate/internal/config"
	"github.com/nix-united/golang-echo-boilerplate/internal/db"
	"github.com/nix-united/golang-echo-boilerplate/internal/slogx"
	"github.com/nix-united/golang-echo-boilerplate/tests/setup"

	"gorm.io/gorm"
)

var gormDB *gorm.DB

func TestMain(m *testing.M) {
	ctx := context.Background()

	shutdown, err := setupMain(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to setup integration tests", "err", err.Error())
		os.Exit(1)
	}

	code := m.Run()

	if err := shutdown(ctx); err != nil {
		slog.ErrorContext(ctx, "Failed to shutdown integration tests", "err", err.Error())
		os.Exit(1)
	}

	os.Exit(code)
}

func setupMain(ctx context.Context) (_ func(context.Context) error, err error) {
	err = slogx.Init(config.LogConfig{
		Application: "integration-tests",
		Level:       "DEBUG",
		AddSource:   true,
	})
	if err != nil {
		return nil, fmt.Errorf("init slog: %w", err)
	}

	shutdownCallbacks := make([]func(context.Context) error, 0)

	shutdown := func(ctx context.Context) error {
		var err error
		for _, callback := range slices.Backward(shutdownCallbacks) {
			err = errors.Join(err, callback(ctx))
		}

		if err != nil {
			return fmt.Errorf("shutdown callbacks: %w", err)
		}

		return nil
	}

	defer func() {
		if err == nil {
			return
		}

		if errShutdown := shutdown(context.Background()); errShutdown != nil {
			err = errors.Join(err, errShutdown)
		}
	}()

	mysqlConfig, mysqlShutdown, err := setup.SetupMySQL(ctx)
	if err != nil {
		return nil, fmt.Errorf("setup mysql: %w", err)
	}

	shutdownCallbacks = append(shutdownCallbacks, mysqlShutdown)

	gormDB, err = db.NewGormDB(config.DBConfig{
		User:     mysqlConfig.User,
		Password: mysqlConfig.Password,
		Name:     mysqlConfig.Name,
		Host:     mysqlConfig.Host,
		Port:     mysqlConfig.ExposedPort,
	})
	if err != nil {
		return nil, fmt.Errorf("new gorm db connection: %w", err)
	}

	dbConnection, err := gormDB.DB()
	if err != nil {
		return nil, fmt.Errorf("get db connection: %w", err)
	}

	if err := setup.MigrateDB(dbConnection); err != nil {
		return nil, fmt.Errorf("run db migrations: %w", err)
	}

	shutdownCallbacks = append(shutdownCallbacks, func(ctx context.Context) error {
		if err := dbConnection.Close(); err != nil {
			return fmt.Errorf("close db connection: %w", err)
		}

		return nil
	})

	return shutdown, nil
}
