package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nix-united/golang-echo-boilerplate/internal/requests"
	"github.com/nix-united/golang-echo-boilerplate/internal/server/handlers"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func newOAuthHandler(t *testing.T) (*echo.Echo, *handlers.OAuthHandler, *MockUserAuthenticator) {
	t.Helper()

	ctrl := gomock.NewController(t)
	userAuthenticator := NewMockUserAuthenticator(ctrl)
	oAuthHandler := handlers.NewOAuthHandler(userAuthenticator)
	engine := echo.New()

	engine.POST("/google-oauth", oAuthHandler.GoogleOAuth)

	return engine, oAuthHandler, userAuthenticator
}

func TestOAuthHandler_GoogleOAuth(t *testing.T) {
	t.Run("It should return an error if received empty request", func(t *testing.T) {
		engine, oAuthHandler, _ := newOAuthHandler(t)

		request := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/google-oauth", http.NoBody)
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		recorder := httptest.NewRecorder()
		c := engine.NewContext(request, recorder)

		err := oAuthHandler.GoogleOAuth(c)
		require.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, recorder.Result().StatusCode)

		wantResponse := `{
			"code": 400,
			"error": "Failed to bind request"
		}`

		assert.JSONEq(t, wantResponse, recorder.Body.String())
	})

	t.Run("It should authorize user", func(t *testing.T) {
		engine, oAuthHandler, userAuthenticator := newOAuthHandler(t)

		oAuthRequest := requests.OAuthRequest{
			Token: "test token",
		}

		buffer := new(bytes.Buffer)
		err := json.NewEncoder(buffer).Encode(oAuthRequest)
		require.NoError(t, err)

		userAuthenticator.
			EXPECT().GoogleOAuth(gomock.Any(), oAuthRequest.Token).
			Return("access-token-123", "refresh-token-456", 3600, nil)

		request := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/google-oauth", buffer)
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		recorder := httptest.NewRecorder()
		c := engine.NewContext(request, recorder)

		err = oAuthHandler.GoogleOAuth(c)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, recorder.Result().StatusCode)

		wantResponse := `{
			  "accessToken": "access-token-123",
              "refreshToken": "refresh-token-456",
              "exp": 3600
		}`

		assert.JSONEq(t, wantResponse, recorder.Body.String())
	})
}
