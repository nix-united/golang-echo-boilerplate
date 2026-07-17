package handlers

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/nix-united/golang-echo-boilerplate/internal/domain"
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
	UpdateByUser(ctx context.Context, request domain.UpdatePostRequest) (*models.Post, error)
	DeleteByUser(ctx context.Context, request domain.DeletePostRequest) error
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
//	@Success		201		{object}	responses.MessageResponse
//	@Failure		400		{object}	responses.ErrorResponse
//	@Security		ApiKeyAuth
//	@Router			/posts [post]
func (p *PostHandlers) CreatePost(c echo.Context) error {
	authClaims, err := getAuthClaims(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, responses.NewErrorResponse("Unauthorized", http.StatusUnauthorized))
	}

	var createPostRequest requests.CreatePostRequest
	if err := c.Bind(&createPostRequest); err != nil {
		return c.JSON(http.StatusBadRequest, responses.NewErrorResponse("Failed to bind request: "+err.Error(), http.StatusBadRequest))
	}

	if err := createPostRequest.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, responses.NewErrorResponse("Required fields are empty", http.StatusBadRequest))
	}

	post := &models.Post{
		Title:   createPostRequest.Title,
		Content: createPostRequest.Content,
		UserID:  authClaims.ID,
	}

	if err := p.postService.Create(c.Request().Context(), post); err != nil {
		return c.JSON(http.StatusBadRequest, responses.NewErrorResponse("Failed to create post: "+err.Error(), http.StatusBadRequest))
	}

	return c.JSON(http.StatusCreated, responses.NewMessageResponse("Post successfully created"))
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
		return c.JSON(http.StatusNotFound, responses.NewErrorResponse("Failed to get all posts: "+err.Error(), http.StatusNotFound))
	}

	response := responses.NewPostResponse(posts)
	return c.JSON(http.StatusOK, response)
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
//	@Success		200		{object}	responses.MessageResponse
//	@Failure		400		{object}	responses.ErrorResponse
//	@Failure		404		{object}	responses.ErrorResponse
//	@Security		ApiKeyAuth
//	@Router			/posts/{id} [put]
func (p *PostHandlers) UpdatePost(c echo.Context) error {
	auth, err := getAuthClaims(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, responses.NewErrorResponse("Unauthorized", http.StatusUnauthorized))
	}

	parsedPostID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, responses.NewErrorResponse("Failed to parse post id: "+err.Error(), http.StatusBadRequest))
	}

	postID, err := safecast.Convert[uint](parsedPostID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, responses.NewErrorResponse("Failed to parse post id: "+err.Error(), http.StatusBadRequest))
	}

	var updatePostRequest requests.UpdatePostRequest
	if err := c.Bind(&updatePostRequest); err != nil {
		return c.JSON(http.StatusBadRequest, responses.NewErrorResponse("Failed to bind request: "+err.Error(), http.StatusBadRequest))
	}

	if err := updatePostRequest.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, responses.NewErrorResponse("Required fields are empty", http.StatusBadRequest))
	}

	_, err = p.postService.UpdateByUser(c.Request().Context(), domain.UpdatePostRequest{
		UserID:  auth.ID,
		PostID:  postID,
		Title:   updatePostRequest.Title,
		Content: updatePostRequest.Content,
	})
	if err != nil {
		switch {
		case errors.Is(err, models.ErrPostNotFound):
			return c.JSON(http.StatusNotFound, responses.NewErrorResponse("Post not found", http.StatusNotFound))
		case errors.Is(err, models.ErrForbidden):
			return c.JSON(http.StatusForbidden, responses.NewErrorResponse("Forbidden", http.StatusForbidden))
		default:
			errorResponse := responses.NewErrorResponse("Failed to update post: "+err.Error(), http.StatusInternalServerError)
			return c.JSON(http.StatusInternalServerError, errorResponse)
		}
	}

	return c.JSON(http.StatusOK, responses.NewMessageResponse("Post successfully updated"))
}

// DeletePost godoc
//
//	@Summary		Delete post
//	@Description	Delete post
//	@ID				posts-delete
//	@Tags			Posts Actions
//	@Param			id	path		int	true	"Post ID"
//	@Success		204	"No Content"
//	@Failure		404	{object}	responses.ErrorResponse
//	@Security		ApiKeyAuth
//	@Router			/posts/{id} [delete]
func (p *PostHandlers) DeletePost(c echo.Context) error {
	auth, err := getAuthClaims(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, responses.NewErrorResponse("Unauthorized", http.StatusUnauthorized))
	}

	parsedID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, responses.NewErrorResponse("Failed to parse post id: "+err.Error(), http.StatusBadRequest))
	}

	postID, err := safecast.Convert[uint](parsedID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, responses.NewErrorResponse("Failed to parse post id: "+err.Error(), http.StatusBadRequest))
	}

	err = p.postService.DeleteByUser(c.Request().Context(), domain.DeletePostRequest{
		UserID: auth.ID,
		PostID: postID,
	})
	if err != nil {
		switch {
		case errors.Is(err, models.ErrPostNotFound):
			return c.JSON(http.StatusNotFound, responses.NewErrorResponse("Post not found", http.StatusNotFound))
		case errors.Is(err, models.ErrForbidden):
			return c.JSON(http.StatusForbidden, responses.NewErrorResponse("Forbidden", http.StatusForbidden))
		default:
			errorResponse := responses.NewErrorResponse("Failed to delete post: "+err.Error(), http.StatusInternalServerError)
			return c.JSON(http.StatusInternalServerError, errorResponse)
		}
	}

	return c.NoContent(http.StatusNoContent)
}
