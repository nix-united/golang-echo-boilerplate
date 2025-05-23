package slogx

import (
	"fmt"
	"io"
	"log/slog"
	"os"
)

type Config struct {
	Application string `env:"LOG_APPLICATION"`

	// File represents path to file where store logs. Used [os.Stdout] if empty.
	File string `env:"LOG_FILE"`

	// One of: "DEBUG", "INFO", "WARN", "ERROR". Default: "DEBUG".
	Level string `env:"LOG_LEVEL"`

	// Add source code position to messages.
	AddSource bool `env:"LOG_ADD_SOURCE"`
}

func Init(config Config) (err error) {
	writer := io.Writer(os.Stdout)
	if config.File != "" {
		writer, err = os.OpenFile(config.File, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
		if err != nil {
			return fmt.Errorf("open file %s: %w", config.File, err)
		}
	}

	level := slog.LevelDebug
	if config.Level != "" {
		if err = level.UnmarshalText([]byte(config.Level)); err != nil {
			return fmt.Errorf("parse log level %s: %w", config.Level, err)
		}
	}

	hostname, err := os.Hostname()
	if err != nil {
		return fmt.Errorf("get hostname: %w", err)
	}

	jsonHandler := slog.NewJSONHandler(writer, &slog.HandlerOptions{AddSource: config.AddSource, Level: level})

	traceHandler := newTraceHandler(jsonHandler)

	logger := slog.New(traceHandler).
		With("application", config.Application).
		With("hostname", hostname)

	slog.SetDefault(logger)

	return nil
}
