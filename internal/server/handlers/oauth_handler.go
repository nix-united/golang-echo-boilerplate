package handlers

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/nix-united/golang-echo-boilerplate/internal/requests"
	"github.com/nix-united/golang-echo-boilerplate/internal/responses"
	"net/http"
)

//go:generate go tool mockgen -source=$GOFILE -destination=oauth_handler_mock_test.go -package=${GOPACKAGE}_test -typed=true

type UserAuthenticator interface {
	GoogleOAuth(ctx context.Context, token string) (string, string, int64, error)
}

type OAuthHandler struct {
	userService UserAuthenticator
}

func NewOAuthHandler(userService UserAuthenticator) *OAuthHandler {
	return &OAuthHandler{userService: userService}
}

// GoogleOAuth godoc
//
//	@Summary		Authenticate user using google provider
//	@Description	Perform user login using google provider
//	@ID				user-auth-google
//	@Tags			User Actions
//	@Accept			json
//	@Produce		json
//	@Param			params	body		requests.OAuthRequest	true	"Google Token"
//	@Success		200		{object}	responses.LoginResponse
//	@Failure		401		{object}	responses.Error
//	@Router			/google-oauth [post]
func (oa *OAuthHandler) GoogleOAuth(c echo.Context) error {
	var oAuthRequest requests.OAuthRequest

	if err := c.Bind(&oAuthRequest); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Failed to bind request")
	}

	if err := oAuthRequest.Validate(); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Required fields are empty or invalid")
	}

	accessToken, refreshToken, exp, err := oa.userService.GoogleOAuth(c.Request().Context(), oAuthRequest.Token)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Failed to authenticate with Google: "+err.Error())
	}

	res := responses.NewLoginResponse(accessToken, refreshToken, exp)
	return responses.Response(c, http.StatusOK, res)
}
