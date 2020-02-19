package builders

import (
	"echo-demo-project/server/models"
)

type PostBuilder struct {
	title   string
	content string
	userId  uint
}

func NewPostBuilder() *PostBuilder {
	return &PostBuilder{}
}

func (postBuilder *PostBuilder) SetTitle(title string) (p *PostBuilder) {
	postBuilder.title = title
	return postBuilder
}
func (postBuilder *PostBuilder) SetContent(content string) (p *PostBuilder) {
	postBuilder.content = content
	return postBuilder
}

func (postBuilder *PostBuilder) SetUserId(userId uint) (p *PostBuilder) {
	postBuilder.userId = userId
	return postBuilder
}

func (postBuilder *PostBuilder) Build() models.Post {
	post := models.Post{
		Title:   postBuilder.title,
		Content: postBuilder.content,
		UserId:  postBuilder.userId,
	}

	return post
}
