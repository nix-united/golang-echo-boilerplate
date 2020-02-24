package handlers

import (
	"echo-demo-project/server/models"
	"echo-demo-project/server/repositories"
	"echo-demo-project/server/requests"
	"echo-demo-project/server/responses"
	"echo-demo-project/server/services"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
)

func (postHandler *PostHandler) UpdatePost() echo.HandlerFunc {
	return func(c echo.Context) error {
		updatePostRequest := new(requests.UpdatePostRequest)
		id, _ := strconv.Atoi(c.Param("id"))

		if err := c.Bind(updatePostRequest); err != nil {
			return err
		}

		if err := c.Validate(updatePostRequest); err != nil {
			return responses.ErrorResponse(c, http.StatusBadRequest, "Required fields are empty")
		}

		post := models.Post{}

		postRepository := repositories.NewPostRepository(postHandler.server.Db)
		postRepository.GetPost(&post, id)

		if post.ID == 0 {
			return responses.ErrorResponse(c, http.StatusNotFound, "Post not found")
		}

		postService := services.NewPostService(postHandler.server.Db)
		postService.Update(&post, updatePostRequest)

		return responses.SuccessResponse(c, "Post successfully update")
	}
}
