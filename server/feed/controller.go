package feed

import (
	"context"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/Jesuloba-world/social-sum/server/database"
)

type createPostSerializer struct {
	Message string `json:"message"`
	Post    *Post  `json:"post"`
}

func getPosts(c *fiber.Ctx) error {
	return c.Status(http.StatusOK).JSON([]Post{{
		Title:   "First Post",
		Content: "This is the first post!",
		// ImageURL: "images/cook.jpg",
		Creator: creator{
			Name: "John Needle",
		},
		// CreatedAt: time.Now(),
	}})
}

func createPost(c *fiber.Ctx) error {
	PostCollection := database.Client.Database("Feed").Collection("Post")

	post := new(Post)
	if err := c.BodyParser(post); err != nil {
		return c.Status(http.StatusBadRequest).SendString(err.Error())
	}

	post.SetTimestamps()
	post.Creator = creator{Name: "Jack Berry"}

	// create post in db
	result, err := PostCollection.InsertOne(context.TODO(), post)
	if err != nil {
		panic(err)
	}

	// Retrieve the inserted document from the database
	insertedPost := new(Post)
	err = PostCollection.FindOne(context.TODO(), bson.M{"_id": result.InsertedID}).Decode(insertedPost)
	if err != nil {
		panic(err)
	}

	return c.Status(http.StatusCreated).JSON(createPostSerializer{Message: "Post created successfully", Post: insertedPost})
}
