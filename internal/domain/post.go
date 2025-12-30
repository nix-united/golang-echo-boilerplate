package domain

type UpdatePostRequest struct {
	// UserID is a user which make request.
	UserID uint

	// PostID is the post to update.
	PostID uint

	Title   string
	Content string
}

type DeletePostRequest struct {
	// UserID is a user which make request.
	UserID uint

	// PostID is the post to update.
	PostID uint
}
