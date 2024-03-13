package feed

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

)

type (
	Post struct {
		ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
		Title     string             `bson:"title" json:"title"`
		Content   string             `bson:"content" json:"content"`
		ImageURL  string             `bson:"imageUrl" json:"imageUrl"`
		CreatorId primitive.ObjectID `bson:"creator" json:"creatorId"`
		Creator   creator            `bson:"_" json:"creator"`
		CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
		UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
	}

	creator struct {
		Name string             `json:"name"`
	}
)

func (p *Post) SetTimestamps() {
	now := time.Now()
	if p.CreatedAt.IsZero() {
		p.CreatedAt = now
	}
	p.UpdatedAt = now
}
