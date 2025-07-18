package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nix-united/golang-echo-boilerplate/internal/requests"
	"github.com/nix-united/golang-echo-boilerplate/internal/responses"
	"github.com/nix-united/golang-echo-boilerplate/internal/server/handlers"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"
)

func newAuthHandler(t *testing.T) (*echo.Echo, *handlers.AuthHandler, *MockauthService) {
	t.Helper()

	ctrl := gomock.NewController(t)
	authService := NewMockauthService(ctrl)
	authHandler := handlers.NewAuthHandler(authService)
	engine := echo.New()

	engine.POST("/login", authHandler.Login)
	engine.POST("/refresh", authHandler.RefreshToken)

	return engine, authHandler, authService
}

func TestAuthHandler_Login(t *testing.T) {
	t.Run("It should return an error when request has no body", func(t *testing.T) {
		engine, authHandler, _ := newAuthHandler(t)

		request := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/login", http.NoBody)
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		recorder := httptest.NewRecorder()
		c := engine.NewContext(request, recorder)

		err := authHandler.Login(c)
		require.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, recorder.Result().StatusCode)

		wantResponse := `{
			"code": 400,
			"error": "Failed to bind request"
		}`

		assert.JSONEq(t, wantResponse, recorder.Body.String())
	})

	t.Run("It should return an error when request is invalid", func(t *testing.T) {
		engine, authHandler, _ := newAuthHandler(t)

		registerRequest := requests.LoginRequest{
			BasicAuth: requests.BasicAuth{
				Email:    "INVALID_EMAIL",
				Password: "some-pass",
			},
		}

		buffer := new(bytes.Buffer)
		err := json.NewEncoder(buffer).Encode(registerRequest)
		require.NoError(t, err)

		request := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/login", buffer)
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		recorder := httptest.NewRecorder()
		c := engine.NewContext(request, recorder)

		err = authHandler.Login(c)
		require.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, recorder.Result().StatusCode)

		wantResponse := `{
			"code": 400,
			"error": "Required fields are empty or not valid"
		}`

		assert.JSONEq(t, wantResponse, recorder.Body.String())
	})

	t.Run("It should authorize user if everything is valid", func(t *testing.T) {
		engine, authHandler, authService := newAuthHandler(t)

		registerRequest := &requests.LoginRequest{
			BasicAuth: requests.BasicAuth{
				Email:    "example@email.com",
				Password: "some-pass",
			},
		}

		buffer := new(bytes.Buffer)
		err := json.NewEncoder(buffer).Encode(registerRequest)
		require.NoError(t, err)

		response := &responses.LoginResponse{
			AccessToken:  "access-token",
			RefreshToken: "refresh-token",
			Exp:          123,
		}

		authService.
			EXPECT().
			GenerateToken(gomock.Any(), registerRequest).
			Return(response, nil)

		request := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/login", buffer)
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		recorder := httptest.NewRecorder()
		c := engine.NewContext(request, recorder)

		err = authHandler.Login(c)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, recorder.Result().StatusCode)

		wantResponse := `{
			"accessToken": "access-token",
			"exp": 123,
			"refreshToken": "refresh-token"
		}`

		assert.JSONEq(t, wantResponse, recorder.Body.String())
	})
}

func TestAuthHandler_RefreshToken(t *testing.T) {
	t.Run("It should return an error when request has no body", func(t *testing.T) {
		engine, registerHandler, _ := newAuthHandler(t)

		request := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/refresh", http.NoBody)
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		recorder := httptest.NewRecorder()
		c := engine.NewContext(request, recorder)

		err := registerHandler.RefreshToken(c)
		require.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, recorder.Result().StatusCode)

		wantResponse := `{
			"code": 400,
			"error": "Failed to bind request"
		}`

		assert.JSONEq(t, wantResponse, recorder.Body.String())
	})

	t.Run("It should refresh token", func(t *testing.T) {
		engine, authHandler, authService := newAuthHandler(t)

		registerRequest := &requests.RefreshRequest{
			Token: "some-token",
		}

		buffer := new(bytes.Buffer)
		err := json.NewEncoder(buffer).Encode(registerRequest)
		require.NoError(t, err)

		response := &responses.LoginResponse{
			AccessToken:  "access-token",
			RefreshToken: "refresh-token",
			Exp:          123,
		}

		authService.
			EXPECT().
			RefreshToken(gomock.Any(), registerRequest).
			Return(response, nil)

		request := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/refresh", buffer)
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		recorder := httptest.NewRecorder()
		c := engine.NewContext(request, recorder)

		err = authHandler.RefreshToken(c)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, recorder.Result().StatusCode)

		wantResponse := `{
			"accessToken": "access-token",
			"exp": 123,
			"refreshToken": "refresh-token"
		}`

		assert.JSONEq(t, wantResponse, recorder.Body.String())
	})
}
