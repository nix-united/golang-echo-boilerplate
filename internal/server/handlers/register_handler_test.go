package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nix-united/golang-echo-boilerplate/internal/models"
	"github.com/nix-united/golang-echo-boilerplate/internal/requests"
	"github.com/nix-united/golang-echo-boilerplate/internal/responses"
	"github.com/nix-united/golang-echo-boilerplate/internal/server/handlers"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestRegisterHandler_Register(t *testing.T) {
	registerRequest := requests.RegisterRequest{
		BasicAuth: requests.BasicAuth{
			Email:    "example@email.com",
			Password: "some-pass",
		},
		Name: "test name",
	}

	testCases := map[string]struct {
		setExpectations func(userRegisterer *MockuserRegisterer)
		request         any
		wantStatus      int
		wantResponse    any
	}{
		"It should return a 400 status code when received empty request": {
			setExpectations: func(userRegisterer *MockuserRegisterer) {},
			request:         map[string]any{},
			wantStatus:      http.StatusBadRequest,
			wantResponse: responses.Error{
				Code:  http.StatusBadRequest,
				Error: "Required fields are empty or invalid",
			},
		},
		"It should return a 400 status code when received invalid request": {
			setExpectations: func(userRegisterer *MockuserRegisterer) {},
			request: requests.RegisterRequest{
				BasicAuth: requests.BasicAuth{
					Email:    "invalid_email",
					Password: "some-pass",
				},
				Name: "test name",
			},
			wantStatus: http.StatusBadRequest,
			wantResponse: responses.Error{
				Code:  http.StatusBadRequest,
				Error: "Required fields are empty or invalid",
			},
		},
		"It should return an error if user exists": {
			setExpectations: func(userRegisterer *MockuserRegisterer) {
				userRegisterer.
					EXPECT().
					GetUserByEmail(gomock.Any(), "example@email.com").
					Return(models.User{Email: "example@email.com"}, nil)
			},
			request:    registerRequest,
			wantStatus: http.StatusConflict,
			wantResponse: responses.Error{
				Code:  http.StatusConflict,
				Error: "User already exists",
			},
		},
		"It should register an user": {
			setExpectations: func(userRegisterer *MockuserRegisterer) {
				userRegisterer.
					EXPECT().
					GetUserByEmail(gomock.Any(), "example@email.com").
					Return(models.User{}, models.ErrUserNotFound)

				userRegisterer.
					EXPECT().
					Register(gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, gotRegisterRequest *requests.RegisterRequest) error {
						assert.Equal(t, &registerRequest, gotRegisterRequest)

						return nil
					})
			},
			request:    registerRequest,
			wantStatus: http.StatusCreated,
			wantResponse: responses.Data{
				Code:    http.StatusCreated,
				Message: "User successfully created",
			},
		},
	}

	for testName, testCase := range testCases {
		t.Run(testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			userRegisterer := NewMockuserRegisterer(ctrl)
			registerHandler := handlers.NewRegisterHandler(userRegisterer)

			testCase.setExpectations(userRegisterer)

			rawRequest, err := json.Marshal(testCase.request)
			require.NoError(t, err)

			request := httptest.NewRequestWithContext(
				t.Context(),
				http.MethodPost,
				"/register",
				bytes.NewBuffer(rawRequest),
			)
			request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			recorder := httptest.NewRecorder()
			c := echo.New().NewContext(request, recorder)

			err = registerHandler.Register(c)
			require.NoError(t, err)

			assert.Equal(t, testCase.wantStatus, recorder.Result().StatusCode)

			wantResponse, err := json.Marshal(testCase.wantResponse)
			require.NoError(t, err)

			assert.JSONEq(t, string(wantResponse), recorder.Body.String())
		})
	}
}
