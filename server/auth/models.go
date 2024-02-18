package auth

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID   `bson:"_id,omitempty" json:"_id"`
	Email     string               `bson:"email" json:"email"`
	Password  string               `bson:"password" json:"-"`
	Name      string               `bson:"name" json:"name"`
	Status    string               `bson:"status" json:"status"`
	Posts     []primitive.ObjectID `bson:"posts" json:"posts"`
	CreatedAt time.Time            `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time            `bson:"updatedAt" json:"updatedAt"`
}

func (u *User) SetTimestamps() {
	now := time.Now()
	if u.CreatedAt.IsZero() {
		u.CreatedAt = now
	}
	u.UpdatedAt = now
}
