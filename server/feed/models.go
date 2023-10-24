package feed

import "time"

type (
	creator struct {
		Name string `json:"name"`
	}

	Post struct {
		Id        string    `json:"_id"`
		Title     string    `json:"title" validate:"required,min=5"`
		Content   string    `json:"content" validate:"required,min=5"`
		ImageUrl  string    `json:"imageUrl"`
		Creator   creator   `json:"creator"`
		CreatedAt time.Time `json:"createdAt"`
	}
)
