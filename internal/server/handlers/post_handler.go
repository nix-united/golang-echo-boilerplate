package handlers

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/nix-united/golang-echo-boilerplate/internal/models"
	"github.com/nix-united/golang-echo-boilerplate/internal/requests"
	"github.com/nix-united/golang-echo-boilerplate/internal/responses"

	safecast "github.com/ccoveille/go-safecast"
	"github.com/labstack/echo/v4"
)

//go:generate go tool mockgen -source=$GOFILE -destination=post_handler_mock_test.go -package=${GOPACKAGE}_test -typed=true

type postService interface {
	Create(ctx context.Context, post *models.Post) error
	GetPosts(ctx context.Context) ([]models.Post, error)
	GetPost(ctx context.Context, id uint) (models.Post, error)
	Update(ctx context.Context, post *models.Post, updatePostRequest requests.UpdatePostRequest) error
	Delete(ctx context.Context, post *models.Post) error
}

type PostHandlers struct {
	postService postService
}

func NewPostHandlers(postService postService) *PostHandlers {
	return &PostHandlers{postService: postService}
}

// CreatePost godoc
//
//	@Summary		Create post
//	@Description	Create post
//	@ID				posts-create
//	@Tags			Posts Actions
//	@Accept			json
//	@Produce		json
//	@Param			params	body		requests.CreatePostRequest	true	"Post title and content"
//	@Success		201		{object}	responses.Data
//	@Failure		400		{object}	responses.Error
//	@Security		ApiKeyAuth
//	@Router			/posts [post]
func (p *PostHandlers) CreatePost(c echo.Context) error {
	authClaims, err := getAuthClaims(c)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
	}

	var createPostRequest requests.CreatePostRequest
	if err := c.Bind(&createPostRequest); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Failed to bind request: "+err.Error())
	}

	if err := createPostRequest.Validate(); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Required fields are empty")
	}

	post := &models.Post{
		Title:   createPostRequest.Title,
		Content: createPostRequest.Content,
		UserID:  authClaims.ID,
	}

	if err := p.postService.Create(c.Request().Context(), post); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Failed to create post: "+err.Error())
	}

	return responses.MessageResponse(c, http.StatusCreated, "Post successfully created")
}

// GetPosts godoc
//
//	@Summary		Get posts
//	@Description	Get the list of all posts
//	@ID				posts-get
//	@Tags			Posts Actions
//	@Produce		json
//	@Success		200	{array}	responses.PostResponse
//	@Security		ApiKeyAuth
//	@Router			/posts [get]
func (p *PostHandlers) GetPosts(c echo.Context) error {
	posts, err := p.postService.GetPosts(c.Request().Context())
	if err != nil {
		return responses.ErrorResponse(c, http.StatusNotFound, "Failed to get all posts: "+err.Error())
	}

	response := responses.NewPostResponse(posts)
	return responses.Response(c, http.StatusOK, response)
}

// UpdatePost godoc
//
//	@Summary		Update post
//	@Description	Update post
//	@ID				posts-update
//	@Tags			Posts Actions
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int							true	"Post ID"
//	@Param			params	body		requests.UpdatePostRequest	true	"Post title and content"
//	@Success		200		{object}	responses.Data
//	@Failure		400		{object}	responses.Error
//	@Failure		404		{object}	responses.Error
//	@Security		ApiKeyAuth
//	@Router			/posts/{id} [put]
func (p *PostHandlers) UpdatePost(c echo.Context) error {
	auth, err := getAuthClaims(c)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
	}

	parsedID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Failed to parse post id: "+err.Error())
	}

	id, err := safecast.ToUint(parsedID)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Failed to parse post id: "+err.Error())
	}

	var updatePostRequest requests.UpdatePostRequest
	if err := c.Bind(&updatePostRequest); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Failed to bind request: "+err.Error())
	}

	if err := updatePostRequest.Validate(); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Required fields are empty")
	}

	post, err := p.postService.GetPost(c.Request().Context(), id)
	if errors.Is(err, models.ErrPostNotFound) {
		return responses.ErrorResponse(c, http.StatusNotFound, "Post not found")
	} else if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, "Failed to find post: "+err.Error())
	}

	if post.UserID != auth.ID {
		return responses.ErrorResponse(c, http.StatusForbidden, "Forbidden")
	}

	if err := p.postService.Update(c.Request().Context(), &post, updatePostRequest); err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, "Failed to update post: "+err.Error())
	}

	return responses.MessageResponse(c, http.StatusOK, "Post successfully updated")
}

// DeletePost godoc
//
//	@Summary		Delete post
//	@Description	Delete post
//	@ID				posts-delete
//	@Tags			Posts Actions
//	@Param			id	path		int	true	"Post ID"
//	@Success		204	{object}	responses.Data
//	@Failure		404	{object}	responses.Error
//	@Security		ApiKeyAuth
//	@Router			/posts/{id} [delete]
func (p *PostHandlers) DeletePost(c echo.Context) error {
	auth, err := getAuthClaims(c)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
	}

	parsedID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Failed to parse post id: "+err.Error())
	}

	id, err := safecast.ToUint(parsedID)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Failed to parse post id: "+err.Error())
	}

	post, err := p.postService.GetPost(c.Request().Context(), id)
	if errors.Is(err, models.ErrPostNotFound) {
		return responses.ErrorResponse(c, http.StatusNotFound, "Post not found")
	} else if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, "Failed to find post: "+err.Error())
	}

	if post.UserID != auth.ID {
		return responses.ErrorResponse(c, http.StatusForbidden, "Forbidden")
	}

	if err := p.postService.Delete(c.Request().Context(), &post); err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete post: "+err.Error())
	}

	return responses.MessageResponse(c, http.StatusNoContent, "Post deleted successfully")
}
