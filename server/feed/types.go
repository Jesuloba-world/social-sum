package feed

import "mime/multipart"

type Error struct {
	Message string `json:"message"`
	Errors  string `json:"error"`
}

type createPostInput struct {
	Title   string                `json:"title"`
	Content string                `json:"content"`
	Image   *multipart.FileHeader `json:"imageUrl"`
}
