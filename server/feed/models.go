package feed

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	creator struct {
		Name string `json:"name"`
	}

	Post struct {
		ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
		Title     string             `bson:"title" json:"title"`
		Content   string             `bson:"content" json:"content"`
		ImageURL  string             `bson:"imageUrl" json:"imageUrl"`
		Creator   creator            `bson:"creator" json:"creator"`
		CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
		UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
	}
)

func (p *Post) SetTimestamps() {
	now := time.Now()
	if p.CreatedAt.IsZero() {
		p.CreatedAt = now
	}
	p.UpdatedAt = now
}
