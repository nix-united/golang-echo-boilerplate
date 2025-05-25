package user_test

import (
	"context"
	"testing"

	"github.com/nix-united/golang-echo-boilerplate/internal/models"
	"github.com/nix-united/golang-echo-boilerplate/internal/requests"
	"github.com/nix-united/golang-echo-boilerplate/internal/services/user"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
)

func TestService_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	userRepository := NewMockuserRepository(ctrl)
	userService := user.NewService(userRepository)

	request := &requests.RegisterRequest{
		BasicAuth: requests.BasicAuth{
			Email:    "example@email.com",
			Password: "some-password",
		},
		Name: "name",
	}

	wantUser := &models.User{
		Email: "example@email.com",
		Name:  "name",
	}

	userRepository.
		EXPECT().
		Create(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, got *models.User) error {
			err := bcrypt.CompareHashAndPassword([]byte(got.Password), []byte("some-password"))
			require.NoError(t, err)

			wantUser.Password = got.Password

			assert.Equal(t, wantUser, got)

			return nil
		})

	err := userService.Register(t.Context(), request)
	require.NoError(t, err)
}

func TestService_GetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	userRepository := NewMockuserRepository(ctrl)
	userService := user.NewService(userRepository)

	wantUser := models.User{
		Email:    "example@email.com",
		Name:     "name",
		Password: "hashed password",
	}

	userRepository.
		EXPECT().
		GetByID(gomock.Any(), uint(123)).
		Return(wantUser, nil)

	gotUser, err := userService.GetByID(t.Context(), 123)
	require.NoError(t, err)

	assert.Equal(t, wantUser, gotUser)
}

func TestService_GetUserByEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	userRepository := NewMockuserRepository(ctrl)
	userService := user.NewService(userRepository)

	wantUser := models.User{
		Email:    "example@gmail.com",
		Name:     "name",
		Password: "hashed password",
	}

	userRepository.
		EXPECT().
		GetUserByEmail(gomock.Any(), "example@gmail.com").
		Return(wantUser, nil)

	gotUser, err := userService.GetUserByEmail(t.Context(), "example@gmail.com")
	require.NoError(t, err)

	assert.Equal(t, wantUser, gotUser)
}
