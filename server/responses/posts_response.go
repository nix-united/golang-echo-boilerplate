package responses

import (
	"echo-demo-project/server/models"
)

type PostResponse struct {
	Title    string
	Content  string
	Username string
	ID       uint
}

func NewPostResponse(posts []models.Post) *[]PostResponse {
	postResponse := make([]PostResponse, 0)

	for i := range posts {
		postResponse = append(postResponse, struct {
			Title    string
			Content  string
			Username string
			ID       uint
		}{
			Title:    posts[i].Title,
			Content:  posts[i].Content,
			Username: posts[i].User.Name,
			ID:       posts[i].ID,
		})
	}

	return &postResponse
}
