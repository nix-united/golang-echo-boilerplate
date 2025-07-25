package auth_test

import (
	"errors"
	"testing"

	"github.com/nix-united/golang-echo-boilerplate/internal/models"
	"github.com/nix-united/golang-echo-boilerplate/internal/requests"
	"github.com/nix-united/golang-echo-boilerplate/internal/responses"
	"github.com/nix-united/golang-echo-boilerplate/internal/services/auth"
	"github.com/nix-united/golang-echo-boilerplate/internal/services/token"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type serviceMocks struct {
	userService  *MockuserService
	tokenService *MocktokenService
}

func newService(t *testing.T) (*auth.Service, serviceMocks) {
	t.Helper()

	ctrl := gomock.NewController(t)
	userService := NewMockuserService(ctrl)
	tokenService := NewMocktokenService(ctrl)
	authService := auth.NewService(userService, tokenService)

	mocks := serviceMocks{
		userService:  userService,
		tokenService: tokenService,
	}

	return authService, mocks
}

func TestService_GenerateToken(t *testing.T) {
	password, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	require.NoError(t, err)

	loginRequest := &requests.LoginRequest{
		BasicAuth: requests.BasicAuth{
			Email:    "example@email.com",
			Password: "password",
		},
	}

	user := models.User{
		Model:    gorm.Model{ID: 1},
		Email:    "example@email.com",
		Name:     "name",
		Password: string(password),
	}

	wantResponse := &responses.LoginResponse{
		AccessToken:  "access-token",
		RefreshToken: "refresh-token",
		Exp:          1000,
	}

	t.Run("It should propagate error from user service when failed to fetch user by email", func(t *testing.T) {
		service, mocks := newService(t)

		userServiceErr := errors.New("error from user service")

		mocks.userService.
			EXPECT().
			GetUserByEmail(gomock.Any(), loginRequest.Email).
			Return(models.User{}, userServiceErr)

		_, err := service.GenerateToken(t.Context(), loginRequest)
		assert.ErrorIs(t, err, userServiceErr)
	})

	t.Run("It should return ErrInvalidPassword error when received invalid password", func(t *testing.T) {
		service, mocks := newService(t)

		loginRequestWithInvalidPassword := *loginRequest
		loginRequestWithInvalidPassword.Password = "invalid-password"

		mocks.userService.
			EXPECT().
			GetUserByEmail(gomock.Any(), loginRequest.Email).
			Return(user, nil)

		_, err := service.GenerateToken(t.Context(), &loginRequestWithInvalidPassword)
		assert.ErrorIs(t, err, models.ErrInvalidPassword)
	})

	t.Run("It should generate token", func(t *testing.T) {
		service, mocks := newService(t)

		mocks.userService.
			EXPECT().
			GetUserByEmail(gomock.Any(), loginRequest.Email).
			Return(user, nil)

		mocks.tokenService.
			EXPECT().
			CreateAccessToken(gomock.Any(), &user).
			Return(wantResponse.AccessToken, wantResponse.Exp, nil)

		mocks.tokenService.
			EXPECT().
			CreateRefreshToken(gomock.Any(), &user).
			Return(wantResponse.RefreshToken, nil)

		response, err := service.GenerateToken(t.Context(), loginRequest)
		require.NoError(t, err)

		assert.Equal(t, wantResponse, response)
	})
}

func TestService_RefreshToken(t *testing.T) {
	refreshRequest := &requests.RefreshRequest{
		Token: "token",
	}

	claims := &token.JwtCustomRefreshClaims{
		ID: 1,
	}

	user := models.User{
		Model:    gorm.Model{ID: 1},
		Email:    "example@email.com",
		Name:     "name",
		Password: "password",
	}

	wantResponse := &responses.LoginResponse{
		AccessToken:  "access-token",
		RefreshToken: "refresh-token",
		Exp:          1000,
	}

	t.Run("It should return ErrInvalidAuthToken when failed to validate token", func(t *testing.T) {
		service, mocks := newService(t)

		mocks.tokenService.
			EXPECT().
			ParseRefreshToken(gomock.Any(), refreshRequest.Token).
			Return(nil, errors.New("error from token service"))

		_, err := service.RefreshToken(t.Context(), refreshRequest)
		assert.ErrorIs(t, err, models.ErrInvalidAuthToken)
	})

	t.Run("It should propagate error from user service when failed to fetch user by ID", func(t *testing.T) {
		service, mocks := newService(t)

		userServiceErr := errors.New("error from user service")

		mocks.tokenService.
			EXPECT().
			ParseRefreshToken(gomock.Any(), refreshRequest.Token).
			Return(claims, nil)

		mocks.userService.
			EXPECT().
			GetByID(gomock.Any(), uint(1)).
			Return(models.User{}, userServiceErr)

		_, err := service.RefreshToken(t.Context(), refreshRequest)
		assert.ErrorIs(t, err, userServiceErr)
	})

	t.Run("It should refresh token", func(t *testing.T) {
		service, mocks := newService(t)

		mocks.tokenService.
			EXPECT().
			ParseRefreshToken(gomock.Any(), refreshRequest.Token).
			Return(claims, nil)

		mocks.userService.
			EXPECT().
			GetByID(gomock.Any(), uint(1)).
			Return(user, nil)

		mocks.tokenService.
			EXPECT().
			CreateAccessToken(gomock.Any(), &user).
			Return(wantResponse.AccessToken, wantResponse.Exp, nil)

		mocks.tokenService.
			EXPECT().
			CreateRefreshToken(gomock.Any(), &user).
			Return(wantResponse.RefreshToken, nil)

		response, err := service.RefreshToken(t.Context(), refreshRequest)
		require.NoError(t, err)

		assert.Equal(t, wantResponse, response)
	})
}
