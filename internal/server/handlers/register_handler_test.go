package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nix-united/golang-echo-boilerplate/internal/models"
	"github.com/nix-united/golang-echo-boilerplate/internal/requests"
	"github.com/nix-united/golang-echo-boilerplate/internal/server/handlers"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func newRegisterHandler(t *testing.T) (*echo.Echo, *handlers.RegisterHandler, *MockuserRegisterer) {
	t.Helper()

	ctrl := gomock.NewController(t)
	userRegisterer := NewMockuserRegisterer(ctrl)
	registerHandler := handlers.NewRegisterHandler(userRegisterer)
	engine := echo.New()

	engine.POST("/register", registerHandler.Register)

	return engine, registerHandler, userRegisterer
}

func TestRegisterHandler_Register(t *testing.T) {
	t.Run("It should return an error if failed to parse request", func(t *testing.T) {
		engine, registerHandler, _ := newRegisterHandler(t)

		request := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/register", http.NoBody)
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		recorder := httptest.NewRecorder()
		c := engine.NewContext(request, recorder)

		err := registerHandler.Register(c)
		require.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, recorder.Result().StatusCode)

		wantResponse := `{
			"code": 400,
			"error": "Failed to bind request"
		}`

		assert.JSONEq(t, wantResponse, recorder.Body.String())
	})

	t.Run("It should return an error if received invalid request", func(t *testing.T) {
		engine, registerHandler, _ := newRegisterHandler(t)

		registerRequest := requests.RegisterRequest{
			BasicAuth: requests.BasicAuth{
				Email:    "invalid_email",
				Password: "some-pass",
			},
			Name: "test name",
		}

		buffer := new(bytes.Buffer)
		err := json.NewEncoder(buffer).Encode(registerRequest)
		require.NoError(t, err)

		request := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/register", buffer)
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		recorder := httptest.NewRecorder()
		c := engine.NewContext(request, recorder)

		err = registerHandler.Register(c)
		require.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, recorder.Result().StatusCode)

		wantResponse := `{
			"code": 400,
			"error": "Required fields are empty or invalid"
		}`

		assert.JSONEq(t, wantResponse, recorder.Body.String())
	})

	t.Run("It should return an error if user exists", func(t *testing.T) {
		engine, registerHandler, userRegisterer := newRegisterHandler(t)

		registerRequest := requests.RegisterRequest{
			BasicAuth: requests.BasicAuth{
				Email:    "example@email.com",
				Password: "some-pass",
			},
			Name: "test name",
		}

		buffer := new(bytes.Buffer)
		err := json.NewEncoder(buffer).Encode(registerRequest)
		require.NoError(t, err)

		userRegisterer.
			EXPECT().
			GetUserByEmail(gomock.Any(), "example@email.com").
			Return(models.User{Email: "example@email.com"}, nil)

		request := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/register", buffer)
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		recorder := httptest.NewRecorder()
		c := engine.NewContext(request, recorder)

		err = registerHandler.Register(c)
		require.NoError(t, err)

		assert.Equal(t, http.StatusConflict, recorder.Result().StatusCode)

		wantResponse := `{
			"code": 409,
			"error": "User already exists"
		}`

		assert.JSONEq(t, wantResponse, recorder.Body.String())
	})

	t.Run("It should register an user", func(t *testing.T) {
		engine, registerHandler, userRegisterer := newRegisterHandler(t)

		registerRequest := requests.RegisterRequest{
			BasicAuth: requests.BasicAuth{
				Email:    "example@email.com",
				Password: "some-pass",
			},
			Name: "test name",
		}

		buffer := new(bytes.Buffer)
		err := json.NewEncoder(buffer).Encode(registerRequest)
		require.NoError(t, err)

		userRegisterer.
			EXPECT().
			GetUserByEmail(gomock.Any(), "example@email.com").
			Return(models.User{}, models.ErrUserNotFound)

		userRegisterer.
			EXPECT().
			Register(gomock.Any(), &registerRequest).
			Return(nil)

		request := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/register", buffer)
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		recorder := httptest.NewRecorder()
		c := engine.NewContext(request, recorder)

		err = registerHandler.Register(c)
		require.NoError(t, err)

		assert.Equal(t, http.StatusCreated, recorder.Result().StatusCode)

		wantResponse := `{
			"code": 201,
			"message": "User successfully created"
		}`

		assert.JSONEq(t, wantResponse, recorder.Body.String())
	})
}
