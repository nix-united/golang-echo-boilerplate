package handlers

import (
	s "echo-demo-project/server"
	"echo-demo-project/server/builders"
	"echo-demo-project/server/models"
	"echo-demo-project/server/repositories"
	"echo-demo-project/server/requests"
	"echo-demo-project/server/responses"
	"echo-demo-project/server/services"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

type PostHandlers struct {
	server *s.Server
}

func NewPostHandlers(server *s.Server) *PostHandlers {
	return &PostHandlers{server: server}
}

// CreatePost godoc
// @Summary Create post
// @Description Create post
// @ID posts-create
// @Tags Posts Actions
// @Accept json
// @Produce json
// @Param params body requests.CreatePostRequest true "Post title and content"
// @Success 200 {string} string "Post successfully created"
// @Failure 400 {string} string "Bad request"
// @Security ApiKeyAuth
// @Router /restricted/posts [post]
func (p *PostHandlers) CreatePost(c echo.Context) error {
	createPostRequest := new(requests.CreatePostRequest)

	if err := c.Bind(createPostRequest); err != nil {
		return err
	}

	if err := c.Validate(createPostRequest); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Required fields are empty")
	}

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*services.JwtCustomClaims)
	id := claims.ID

	post := builders.NewPostBuilder().
		SetTitle(createPostRequest.Title).
		SetContent(createPostRequest.Content).
		SetUserID(id).
		Build()
	postService := services.NewPostService(p.server.Db)
	postService.Create(&post)

	return responses.SuccessResponse(c, "Post successfully created")
}

func (p *PostHandlers) DeletePost(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	post := models.Post{}

	postRepository := repositories.NewPostRepository(p.server.Db)
	postRepository.GetPost(&post, id)

	if post.ID == 0 {
		return responses.ErrorResponse(c, http.StatusNotFound, "Post not found")
	}

	postService := services.NewPostService(p.server.Db)
	postService.Delete(&post)

	return responses.SuccessResponse(c, "Post delete successfully")
}

func (p *PostHandlers) GetPosts(c echo.Context) error {
	var posts []models.Post

	postRepository := repositories.NewPostRepository(p.server.Db)
	postRepository.GetPosts(&posts)

	for i := 0; i < len(posts); i++ {
		p.server.Db.Model(&posts[i]).Related(&posts[i].User)
	}

	response := responses.NewPostResponse(posts)
	return responses.SuccessResponse(c, response)
}

func (p *PostHandlers) UpdatePost(c echo.Context) error {
	updatePostRequest := new(requests.UpdatePostRequest)
	id, _ := strconv.Atoi(c.Param("id"))

	if err := c.Bind(updatePostRequest); err != nil {
		return err
	}

	if err := c.Validate(updatePostRequest); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Required fields are empty")
	}

	post := models.Post{}

	postRepository := repositories.NewPostRepository(p.server.Db)
	postRepository.GetPost(&post, id)

	if post.ID == 0 {
		return responses.ErrorResponse(c, http.StatusNotFound, "Post not found")
	}

	postService := services.NewPostService(p.server.Db)
	postService.Update(&post, updatePostRequest)

	return responses.SuccessResponse(c, "Post successfully update")
}
