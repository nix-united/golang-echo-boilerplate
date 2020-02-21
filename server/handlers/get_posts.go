package handlers

import (
	"echo-demo-project/server/models"
	"echo-demo-project/server/repositories"
	"echo-demo-project/server/responses"
	"github.com/labstack/echo"
)

func (postHandler *PostHandler) GetPosts() echo.HandlerFunc {
	return func(c echo.Context) error {
		posts := []models.Post{}

		postRepository := repositories.NewPostRepository(postHandler.server.Db)
		postRepository.GetPosts(&posts)


		for i := 0; i < len(posts); i ++ {
			postHandler.server.Db.Model(&posts[i]).Related(&posts[i].User)
		}

		response := responses.NewPostResponse(posts)
		return responses.SuccessResponse(c, response)
	}
}

