package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

// requestDebugger is a logging middleware that logs request and response bodies with DEBUG level logs.
// Warning: Do not use this middleware with endpoints containing sensitive information.
//
// Based on the Content-Type, it determines how the body will be formatted.
// If the content type is application/json, the body will be logged as JSON; otherwise, it will be logged as a string.
type requestDebugger struct{}

func NewRequestDebugger() echo.MiddlewareFunc {
	return (&requestDebugger{}).handle
}

func (d *requestDebugger) handle(next echo.HandlerFunc) echo.HandlerFunc {
	if !slog.Default().Enabled(context.Background(), slog.LevelDebug) {
		return next
	}

	return func(c echo.Context) error {
		requestBody, err := d.getRequestBody(c)
		if err != nil {
			return fmt.Errorf("get request body for logging: %w", err)
		}

		responseBodyGetter := d.getResponseBodyGetter(c)

		errNext := next(c)

		var attrs []any
		if requestBody != nil {
			attrs = append(attrs, "request_body", requestBody)
		}

		responseBody := responseBodyGetter(c)
		if responseBody != nil {
			attrs = append(attrs, "response_body", responseBody)
		}

		message := "Request/response data"
		if len(attrs) == 0 {
			message = "Request/response withot any data"
		}

		slog.DebugContext(c.Request().Context(), message, attrs...)

		if errNext != nil {
			return fmt.Errorf("handle request with request debugger: %w", errNext)
		}

		return nil
	}
}

func (d *requestDebugger) getRequestBody(c echo.Context) (any, error) {
	if c.Request().Body == nil {
		return nil, nil
	}

	rawRequestBody, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return nil, fmt.Errorf("read request body: %w", err)
	}

	request := c.Request()
	request.Body = io.NopCloser(bytes.NewReader(rawRequestBody))
	c.SetRequest(request)

	if strings.HasPrefix(request.Header.Get(echo.HeaderContentType), echo.MIMEApplicationJSON) {
		return json.RawMessage(rawRequestBody), nil
	}

	return string(rawRequestBody), nil
}

func (d *requestDebugger) getResponseBodyGetter(c echo.Context) func(echo.Context) any {
	response := c.Response()
	storer := newResponseStorer(response.Writer)
	response.Writer = storer
	c.SetResponse(response)

	return func(c echo.Context) any {
		if storer.storedResponse == nil {
			return nil
		}

		if strings.HasPrefix(c.Response().Header().Get(echo.HeaderContentType), echo.MIMEApplicationJSON) {
			return json.RawMessage(storer.storedResponse)
		}

		return string(storer.storedResponse)
	}
}

// responseStorer stores the written response by the handler into its field.
// This is used to automate response logging.
type responseStorer struct {
	http.ResponseWriter
	storedResponse []byte
}

func newResponseStorer(writer http.ResponseWriter) *responseStorer {
	return &responseStorer{ResponseWriter: writer}
}

func (s *responseStorer) Write(response []byte) (int, error) {
	s.storedResponse = response

	n, err := s.ResponseWriter.Write(response)
	if err != nil {
		return n, fmt.Errorf("write response with storer: %w", err)
	}

	return n, nil
}
