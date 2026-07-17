package handlers

import (
	"context"
	"net/http"

	"github.com/nix-united/golang-echo-boilerplate/internal/requests"
	"github.com/nix-united/golang-echo-boilerplate/internal/responses"

	"github.com/labstack/echo/v4"
)

//go:generate go tool mockgen -source=$GOFILE -destination=oauth_handler_mock_test.go -package=${GOPACKAGE}_test -typed=true

type userAuthenticator interface {
	GoogleOAuth(ctx context.Context, token string) (accessToken string, refreshToken string, exp int64, err error)
}

type OAuthHandler struct {
	userService userAuthenticator
}

func NewOAuthHandler(userService userAuthenticator) *OAuthHandler {
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
//	@Failure		401		{object}	responses.ErrorResponse
//	@Router			/google-oauth [post]
func (oa *OAuthHandler) GoogleOAuth(c echo.Context) error {
	var oAuthRequest requests.OAuthRequest

	if err := c.Bind(&oAuthRequest); err != nil {
		return c.JSON(http.StatusBadRequest, responses.NewErrorResponse("Failed to bind request", http.StatusBadRequest))
	}

	if err := oAuthRequest.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, responses.NewErrorResponse("Required fields are empty or invalid", http.StatusBadRequest))
	}

	accessToken, refreshToken, exp, err := oa.userService.GoogleOAuth(c.Request().Context(), oAuthRequest.Token)
	if err != nil {
		errorResponse := responses.NewErrorResponse("Failed to authenticate with Google: "+err.Error(), http.StatusBadRequest)
		return c.JSON(http.StatusBadRequest, errorResponse)
	}

	res := responses.NewLoginResponse(accessToken, refreshToken, exp)
	return c.JSON(http.StatusOK, res)
}
