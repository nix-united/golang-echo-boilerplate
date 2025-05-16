package integration

import (
	"testing"

	"github.com/nix-united/golang-echo-boilerplate/internal/models"
	"github.com/nix-united/golang-echo-boilerplate/internal/repositories"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostRepository(t *testing.T) {
	userRepository := repositories.NewUserRepository(gormDB)
	postRepository := repositories.NewPostRepository(gormDB)

	userToCreate := &models.User{
		Email:    "example_user_email@email.com",
		Name:     "some-user-with-posts",
		Password: "some-user-with-posts-password",
	}

	err := gormDB.Create(userToCreate).Error
	require.NoError(t, err)

	user, err := userRepository.GetUserByEmail("example_user_email@email.com")
	require.NotNil(t, user)

	postToCreate := &models.Post{
		Title:   "Post title",
		Content: "Post content",
		UserID:  user.ID,
	}

	postRepository.Create(postToCreate)

	gotPosts, err := postRepository.GetPosts()
	require.NoError(t, err)
	require.Len(t, gotPosts, 1)

	wantPost := postToCreate
	wantPost.User = user
	wantPost.CreatedAt = gotPosts[0].CreatedAt
	wantPost.UpdatedAt = gotPosts[0].UpdatedAt

	assert.Equal(t, *postToCreate, gotPosts[0])
}
