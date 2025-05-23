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

	user, err := userRepository.GetUserByEmail(t.Context(), "example_user_email@email.com")
	require.NotNil(t, user)

	newPost := &models.Post{
		Title:   "Post title",
		Content: "Post content",
		UserID:  user.ID,
	}

	t.Run("It should create a post", func(t *testing.T) {
		err := postRepository.Create(t.Context(), newPost)
		require.NoError(t, err)
		assert.NotZero(t, newPost.ID)
	})

	t.Run("It should fetch created post", func(t *testing.T) {
		gotPost, err := postRepository.GetPost(t.Context(), newPost.ID)
		require.NoError(t, err)

		newPost.CreatedAt = gotPost.CreatedAt
		newPost.UpdatedAt = gotPost.UpdatedAt

		assert.Equal(t, *newPost, gotPost)
	})

	t.Run("It should return an error if post not found", func(t *testing.T) {
		_, err := postRepository.GetPost(t.Context(), 999)
		assert.ErrorIs(t, err, models.ErrPostNotFound)
	})

	t.Run("It should fetch all posts", func(t *testing.T) {
		posts, err := postRepository.GetPosts(t.Context())
		require.NoError(t, err)
		require.Len(t, posts, 1)
		assert.Equal(t, *newPost, posts[0])
	})

	t.Run("It should update post", func(t *testing.T) {
		newPost.Title = "New post title"
		newPost.Content = "New post content"
		err := postRepository.Update(t.Context(), newPost)
		require.NoError(t, err)

		gotPost, err := postRepository.GetPost(t.Context(), newPost.ID)
		require.NoError(t, err)

		newPost.UpdatedAt = gotPost.UpdatedAt

		assert.Equal(t, *newPost, gotPost)
	})

	t.Run("It should delete post", func(t *testing.T) {
		id := newPost.ID

		err := postRepository.Delete(t.Context(), newPost)
		require.NoError(t, err)

		_, err = postRepository.GetPost(t.Context(), id)
		assert.ErrorIs(t, err, models.ErrPostNotFound)
	})
}
