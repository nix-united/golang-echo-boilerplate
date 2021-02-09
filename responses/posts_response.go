package responses

import (
	"echo-demo-project/models"
)

type PostResponse struct {
	Title    string `json:"title" example:"Echo"`
	Content  string `json:"content" example:"Echo is nice!"`
	Username string `json:"username" example:"John Doe"`
	ID       uint   `json:"id" example:"1"`
}

func NewPostResponse(posts []models.Post) *[]PostResponse {
	postResponse := make([]PostResponse, 0)

	for i := range posts {
		postResponse = append(postResponse, PostResponse{
			Title:    posts[i].Title,
			Content:  posts[i].Content,
			Username: posts[i].User.Name,
			ID:       posts[i].ID,
		})
	}

	return &postResponse
}
