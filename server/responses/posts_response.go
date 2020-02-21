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

	for _, post := range posts {
		postResponse = append(postResponse, struct {
			Title    string
			Content  string
			Username string
			ID       uint
		}{Title: post.Title, Content: post.Content, Username: post.User.Name, ID: post.ID})
	}

	return &postResponse
}
