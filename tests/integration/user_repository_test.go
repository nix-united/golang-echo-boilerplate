package integration

import (
	"testing"

	"github.com/nix-united/golang-echo-boilerplate/internal/models"
	"github.com/nix-united/golang-echo-boilerplate/internal/repositories"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserRepository(t *testing.T) {
	userRepository := repositories.NewUserRepository(gormDB)

	newUser := &models.User{
		Email:    "test_user_repository@email.com",
		Name:     "test_user_repository",
		Password: "test_user_repository",
	}

	t.Run("It should create an user", func(t *testing.T) {
		err := userRepository.Create(t.Context(), newUser)
		require.NoError(t, err)
		assert.NotZero(t, newUser.ID)
	})

	t.Run("It should fetch created user", func(t *testing.T) {
		gotUser, err := userRepository.GetByID(t.Context(), newUser.ID)
		require.NoError(t, err)

		newUser.CreatedAt = gotUser.CreatedAt
		newUser.UpdatedAt = gotUser.UpdatedAt

		assert.Equal(t, *newUser, gotUser)
	})

	t.Run("It should return an error if user with such ID not found", func(t *testing.T) {
		_, err := userRepository.GetByID(t.Context(), 999)
		assert.ErrorIs(t, err, models.ErrUserNotFound)
	})

	t.Run("It should fetch user by email", func(t *testing.T) {
		gotUser, err := userRepository.GetUserByEmail(t.Context(), newUser.Email)
		require.NoError(t, err)

		newUser.CreatedAt = gotUser.CreatedAt
		newUser.UpdatedAt = gotUser.UpdatedAt

		assert.Equal(t, *newUser, gotUser)
	})

	t.Run("It should return an error if user with such email not found", func(t *testing.T) {
		_, err := userRepository.GetUserByEmail(t.Context(), "unknown_email@gmail.com")
		assert.ErrorIs(t, err, models.ErrUserNotFound)
	})
}
