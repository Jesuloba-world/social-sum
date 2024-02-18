package feed

import (
	"encoding/json"
	"fmt"
	"mime/multipart"

)

type Error struct {
	Message string `json:"message"`
	Errors  string `json:"error"`
}

type createPostInput struct {
	Title   string     `json:"title" validate:"required,min=5"`
	Content string     `json:"content" validate:"required,min=5"`
	Image   imageField `json:"imageUrl" validate:"required"`
}

type imageField struct {
	File *multipart.FileHeader
	URL  string
}

func (i *imageField) UnmarshalJSON(data []byte) error {
	// Try to unmarshal as a string
	var url string
	if err := json.Unmarshal(data, &url); err == nil {
		i.URL = url
		return nil
	}

	// If it's not a string, try to unmarshal as a file
	var file multipart.FileHeader
	if err := json.Unmarshal(data, &file); err == nil {
		i.File = &file
		return nil
	}

	return fmt.Errorf("invalid image field")
}
