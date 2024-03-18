package graph

import (
	"context"
	"fmt"
	"time"

	"github.com/Jesuloba-world/social-sum/server/auth"
	"github.com/Jesuloba-world/social-sum/server/feed"
	"github.com/Jesuloba-world/social-sum/server/graph/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserResolver struct{ *Resolver }

func (r *UserResolver) Posts(ctx context.Context, user *model.User) ([]*model.Post, error) {
	var posts []*model.Post

	postCollection := r.DB.Database("Feed").Collection("Post")
	userCollection := r.DB.Database("Auth").Collection("User")

	postObjectIds := make([]primitive.ObjectID, len(user.Posts))
	for i, post := range user.Posts {
		objectId, err := primitive.ObjectIDFromHex(post.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to convert post Id to object Id: %s", err.Error())
		}
		postObjectIds[i] = objectId
	}

	// find the post using object ids
	filter := bson.M{"_id": bson.M{"$in": postObjectIds}}
	cursor, err := postCollection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch posts: %s", err.Error())
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var dbpost feed.Post
		err := cursor.Decode(&dbpost)
		if err != nil {
			return nil, fmt.Errorf("failed to decode post: %s", err.Error())
		}

		post := &model.Post{
			ID:        dbpost.ID.Hex(),
			Title:     dbpost.Title,
			Content:   dbpost.Content,
			ImageURL:  dbpost.ImageURL,
			CreatedAt: dbpost.CreatedAt.Format(time.RFC3339),
			UpdatedAt: dbpost.UpdatedAt.Format(time.RFC3339),
		}

		var dbCreator auth.User
		filter := bson.M{"_id": dbpost.CreatorId}
		err = userCollection.FindOne(ctx, filter).Decode(&dbCreator)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch creator: %s", err.Error())
		}

		creator := &model.User{
			ID:     dbCreator.ID.Hex(),
			Email:  dbCreator.Email,
			Name:   dbCreator.Name,
			Status: dbCreator.Status,
		}

		post.Creator = creator

		posts = append(posts, post)
	}

	return posts, nil
}
