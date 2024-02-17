package feed

import "mime/multipart"

type Error struct {
	Message string `json:"message"`
	Errors  string `json:"error"`
}

type createPostInput struct {
	Title   string                `json:"title" validate:"required,min=5"`
	Content string                `json:"content" validate:"required,min=5"`
	Image   *multipart.FileHeader `json:"imageUrl" validate:"required"`
}
