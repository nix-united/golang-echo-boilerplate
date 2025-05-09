package handlers

import (
	"net/http"
	"strconv"

	"github.com/nix-united/golang-echo-boilerplate/internal/models"
	"github.com/nix-united/golang-echo-boilerplate/internal/requests"
	"github.com/nix-united/golang-echo-boilerplate/internal/responses"
	"github.com/nix-united/golang-echo-boilerplate/internal/services/token"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type postService interface {
	Create(post *models.Post) error
	GetPosts() ([]models.Post, error)
	GetPost(id int) (models.Post, error)
	Update(post *models.Post, updatePostRequest requests.UpdatePostRequest) error
	Delete(post *models.Post) error
}

type PostHandlers struct {
	postService postService
}

func NewPostHandlers(postService postService) PostHandlers {
	return PostHandlers{postService: postService}
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
	var createPostRequest requests.CreatePostRequest
	if err := c.Bind(&createPostRequest); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Failed to bind request: "+err.Error())
	}

	if err := createPostRequest.Validate(); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Required fields are empty")
	}

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*token.JwtCustomClaims)
	id := claims.ID

	post := models.Post{
		Title:   createPostRequest.Title,
		Content: createPostRequest.Content,
		UserID:  id,
	}

	if err := p.postService.Create(&post); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Failed to create post: "+err.Error())
	}

	return responses.MessageResponse(c, http.StatusCreated, "Post successfully created")
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
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Failed to parse post id: "+err.Error())
	}

	post, err := p.postService.GetPost(id)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusNotFound, "Post not found")
	}

	if err := p.postService.Delete(&post); err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete post: "+err.Error())
	}

	return responses.MessageResponse(c, http.StatusNoContent, "Post deleted successfully")
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
	posts, err := p.postService.GetPosts()
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
	id, err := strconv.Atoi(c.Param("id"))
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

	post, err := p.postService.GetPost(id)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusNotFound, "Post not found")
	}

	if err := p.postService.Update(&post, updatePostRequest); err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, "Failed to update post: "+err.Error())
	}

	return responses.MessageResponse(c, http.StatusOK, "Post successfully updated")
}
