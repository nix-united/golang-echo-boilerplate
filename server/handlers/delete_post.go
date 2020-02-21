package handlers

import (
	"echo-demo-project/server/models"
	"echo-demo-project/server/repositories"
	"echo-demo-project/server/responses"
	"echo-demo-project/server/services"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
)

func (postHandler *PostHandler) DeletePost() echo.HandlerFunc {
	return func(c echo.Context) error {
		id, _ := strconv.Atoi(c.Param("id"))

		post := models.Post{}

		postRepository := repositories.NewPostRepository(postHandler.server.Db)
		postRepository.GetPost(&post, id)

		if post.ID == 0 {
			return responses.ErrorResponse(c, http.StatusNotFound, "Post not found")
		}

		postService := services.NewPostService(postHandler.server.Db)
		postService.Delete(&post)

		return responses.SuccessResponse(c, "Post delete successfully")
	}
}
