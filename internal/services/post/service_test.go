package post_test

import (
	"testing"

	"github.com/nix-united/golang-echo-boilerplate/internal/domain"
	"github.com/nix-united/golang-echo-boilerplate/internal/models"
	"github.com/nix-united/golang-echo-boilerplate/internal/services/post"
	"gorm.io/gorm"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestService_Create(t *testing.T) {
	newPost := &models.Post{
		Title:   "title",
		Content: "conent",
		UserID:  111,
	}

	ctrl := gomock.NewController(t)
	postRepository := NewMockpostRepository(ctrl)
	postService := post.NewService(postRepository)

	postRepository.
		EXPECT().
		Create(gomock.Any(), newPost).
		Return(nil)

	err := postService.Create(t.Context(), newPost)
	require.NoError(t, err)
}

func TestService_GetPosts(t *testing.T) {
	wantPosts := []models.Post{{
		Title:   "title",
		Content: "conent",
		UserID:  111,
	}}

	ctrl := gomock.NewController(t)
	postRepository := NewMockpostRepository(ctrl)
	postService := post.NewService(postRepository)

	postRepository.
		EXPECT().
		GetPosts(gomock.Any()).
		Return(wantPosts, nil)

	gotPosts, err := postService.GetPosts(t.Context())
	require.NoError(t, err)

	assert.Equal(t, wantPosts, gotPosts)
}

func TestService_GetPost(t *testing.T) {
	wantPost := models.Post{
		Title:   "title",
		Content: "conent",
		UserID:  111,
	}

	ctrl := gomock.NewController(t)
	postRepository := NewMockpostRepository(ctrl)
	postService := post.NewService(postRepository)

	postRepository.
		EXPECT().
		GetPost(gomock.Any(), uint(123)).
		Return(wantPost, nil)

	gotPost, err := postService.GetPost(t.Context(), 123)
	require.NoError(t, err)

	assert.Equal(t, wantPost, gotPost)
}

func TestService_UpdateByUser(t *testing.T) {
	oldPost := models.Post{
		Model:   gorm.Model{ID: 222},
		Title:   "title",
		Content: "conent",
		UserID:  111,
	}

	wantPost := &models.Post{
		Model:   gorm.Model{ID: 222},
		Title:   "new title",
		Content: "new content",
		UserID:  111,
	}

	request := domain.UpdatePostRequest{
		UserID:  111,
		PostID:  222,
		Title:   "new title",
		Content: "new content",
	}

	ctrl := gomock.NewController(t)
	postRepository := NewMockpostRepository(ctrl)
	postService := post.NewService(postRepository)

	postRepository.
		EXPECT().
		GetPost(gomock.Any(), request.PostID).
		Return(oldPost, nil)

	postRepository.
		EXPECT().
		Update(gomock.Any(), wantPost).
		Return(nil)

	newPost, err := postService.UpdateByUser(t.Context(), request)
	require.NoError(t, err)

	assert.Equal(t, wantPost, newPost)
}

func TestService_Delete(t *testing.T) {
	wantPost := &models.Post{
		Title:   "new title",
		Content: "new content",
		UserID:  111,
	}

	ctrl := gomock.NewController(t)
	postRepository := NewMockpostRepository(ctrl)
	postService := post.NewService(postRepository)

	postRepository.
		EXPECT().
		Delete(gomock.Any(), wantPost).
		Return(nil)

	err := postService.Delete(t.Context(), wantPost)
	require.NoError(t, err)
}
