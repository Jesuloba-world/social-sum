package feed

import "time"

type creator struct {
	Name string `json:"name"`
}

type Post struct {
	Id        string    `json:"_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	ImageUrl  string    `json:"imageUrl"`
	Creator   creator   `json:"creator"`
	CreatedAt time.Time `json:"createdAt"`
}
