package handlers_test

import (
	"bytes"
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
	gomock "go.uber.org/mock/gomock"
)

func newAuthHandler(t *testing.T) (*handlers.AuthHandler, *MockauthService) {
	t.Helper()

	ctrl := gomock.NewController(t)
	authService := NewMockauthService(ctrl)
	authHandler := handlers.NewAuthHandler(authService)

	return authHandler, authService
}

func TestAuthHandler_Login(t *testing.T) {
	invalidRequest := &requests.LoginRequest{
		BasicAuth: requests.BasicAuth{
			Email:    "INVALID_EMAIL",
			Password: "some-pass",
		},
	}

	request := &requests.LoginRequest{
		BasicAuth: requests.BasicAuth{
			Email:    "example@email.com",
			Password: "some-pass",
		},
	}

	response := &responses.LoginResponse{
		AccessToken:  "access-token",
		RefreshToken: "refresh-token",
		Exp:          123,
	}

	testCases := map[string]struct {
		setExpectations func(authService *MockauthService)
		request         any
		wantStatus      int
		wantResponse    any
	}{
		"It should return 400 status code when request is invalid": {
			setExpectations: func(authService *MockauthService) {},
			request:         invalidRequest,
			wantStatus:      http.StatusBadRequest,
			wantResponse: responses.Error{
				Code:  http.StatusBadRequest,
				Error: "Required fields are empty or not valid",
			},
		},
		"It should return 401 status code when user does not exist": {
			setExpectations: func(authService *MockauthService) {
				authService.
					EXPECT().
					GenerateToken(gomock.Any(), request).
					Return(nil, models.ErrUserNotFound)
			},
			request:    request,
			wantStatus: http.StatusUnauthorized,
			wantResponse: responses.Error{
				Code:  http.StatusUnauthorized,
				Error: "Invalid credentials",
			},
		},
		"It should return 401 status code when passoword is invalid": {
			setExpectations: func(authService *MockauthService) {
				authService.
					EXPECT().
					GenerateToken(gomock.Any(), request).
					Return(nil, models.ErrInvalidPassword)
			},
			request:    request,
			wantStatus: http.StatusUnauthorized,
			wantResponse: responses.Error{
				Code:  http.StatusUnauthorized,
				Error: "Invalid credentials",
			},
		},
		"It should authorize user": {
			setExpectations: func(authService *MockauthService) {
				authService.
					EXPECT().
					GenerateToken(gomock.Any(), request).
					Return(response, nil)
			},
			request:      request,
			wantStatus:   http.StatusOK,
			wantResponse: response,
		},
	}

	for testName, testCase := range testCases {
		t.Run(testName, func(t *testing.T) {
			authHandler, authService := newAuthHandler(t)

			testCase.setExpectations(authService)

			rawRequest, err := json.Marshal(testCase.request)
			require.NoError(t, err)

			request := httptest.NewRequestWithContext(
				t.Context(),
				http.MethodPost,
				"/login",
				bytes.NewBuffer(rawRequest),
			)
			request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			recorder := httptest.NewRecorder()
			c := echo.New().NewContext(request, recorder)

			err = authHandler.Login(c)
			require.NoError(t, err)

			assert.Equal(t, testCase.wantStatus, recorder.Result().StatusCode)

			wantResponse, err := json.Marshal(testCase.wantResponse)
			require.NoError(t, err)

			assert.JSONEq(t, string(wantResponse), recorder.Body.String())
		})
	}
}

func TestAuthHandler_RefreshToken(t *testing.T) {
	request := &requests.RefreshRequest{
		Token: "some-token",
	}

	rawRequest, err := json.Marshal(request)
	require.NoError(t, err)

	response := &responses.LoginResponse{
		AccessToken:  "access-token",
		RefreshToken: "refresh-token",
		Exp:          123,
	}

	testCases := map[string]struct {
		setExpectations func(authService *MockauthService)
		wantStatus      int
		wantResponse    any
	}{
		"It should respond with a 401 status code when the user does not exist": {
			setExpectations: func(authService *MockauthService) {
				authService.
					EXPECT().
					RefreshToken(gomock.Any(), request).
					Return(nil, models.ErrUserNotFound)
			},
			wantStatus: http.StatusUnauthorized,
			wantResponse: responses.Error{
				Code:  http.StatusUnauthorized,
				Error: "Unauthorized",
			},
		},
		"It should respond with a 401 status code when token is invalid": {
			setExpectations: func(authService *MockauthService) {
				authService.
					EXPECT().
					RefreshToken(gomock.Any(), request).
					Return(nil, models.ErrInvalidAuthToken)
			},
			wantStatus: http.StatusUnauthorized,
			wantResponse: responses.Error{
				Code:  http.StatusUnauthorized,
				Error: "Unauthorized",
			},
		},
		"It should refresh token": {
			setExpectations: func(authService *MockauthService) {
				authService.
					EXPECT().
					RefreshToken(gomock.Any(), request).
					Return(response, nil)
			},
			wantStatus: http.StatusOK,
			wantResponse: responses.LoginResponse{
				AccessToken:  "access-token",
				RefreshToken: "refresh-token",
				Exp:          123,
			},
		},
	}

	for testName, testCase := range testCases {
		t.Run(testName, func(t *testing.T) {
			authHandler, authService := newAuthHandler(t)

			testCase.setExpectations(authService)

			request := httptest.NewRequestWithContext(
				t.Context(),
				http.MethodPost,
				"/refresh",
				bytes.NewBuffer(rawRequest),
			)
			request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			recorder := httptest.NewRecorder()
			c := echo.New().NewContext(request, recorder)

			err = authHandler.RefreshToken(c)
			require.NoError(t, err)

			assert.Equal(t, testCase.wantStatus, recorder.Result().StatusCode)

			wantResponse, err := json.Marshal(testCase.wantResponse)
			require.NoError(t, err)

			assert.JSONEq(t, string(wantResponse), recorder.Body.String())
		})
	}
}
