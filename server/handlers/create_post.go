package handlers

import (
	"echo-demo-project/server/builders"
	"echo-demo-project/server/requests"
	"echo-demo-project/server/responses"
	"echo-demo-project/server/services"
	"github.com/labstack/echo"
	"net/http"
)

func (postHandler *PostHandler) CreatePost() echo.HandlerFunc {
	return func(c echo.Context) error {
		createPostRequest := new(requests.CreatePostRequest)

		if err := c.Bind(createPostRequest); err != nil {
			return err
		}

		if err := c.Validate(createPostRequest); err != nil {
			return responses.ErrorResponse(c, http.StatusBadRequest, "Required fields are empty")
		}

		post := builders.NewPostBuilder().
			SetTitle(createPostRequest.Title).
			SetContent(createPostRequest.Content).
			SetUserId(1).
			Build()
		postService := services.NewPostService(postHandler.server.Db)
		postService.Create(&post)

		return responses.SuccessResponse(c, "Post successfully create")
	}
}
