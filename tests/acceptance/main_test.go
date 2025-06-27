package acceptance

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"testing"

	"github.com/nix-united/golang-echo-boilerplate/tests/setup"
)

var applicationURL *url.URL

func TestMain(m *testing.M) {
	ctx := context.Background()

	shutdown, err := setupMain(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to setup acceptance tests", "err", err.Error())
		os.Exit(1)
	}

	code := m.Run()

	if err := shutdown(ctx); err != nil {
		slog.ErrorContext(ctx, "Failed to shutdown acceptance tests", "err", err.Error())
		os.Exit(1)
	}

	os.Exit(code)
}

func setupMain(ctx context.Context) (_ func(ctx context.Context) error, err error) {
	url, shutdown, err := setup.SetupApplication(ctx)
	if err != nil {
		return nil, fmt.Errorf("setup application: %w", err)
	}

	applicationURL = url

	return shutdown, nil
}
