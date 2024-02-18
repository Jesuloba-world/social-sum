package feed

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/Jesuloba-world/social-sum/server/database"
)

type postSerializer struct {
	Message string `json:"message"`
	Post    *Post  `json:"post"`
}

type allPostSerializer struct {
	Message string `json:"message"`
	Posts   []Post `json:"posts"`
}

func getPosts(c *fiber.Ctx) error {
	PostCollection := database.Client.Database("Feed").Collection("Post")

	var posts []Post

	// Find all documents in the collection
	cursor, err := PostCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}
	defer cursor.Close(context.TODO())

	// Iterate over the cursor and decode each document into a Post struct
	for cursor.Next(context.TODO()) {
		var post Post
		if err := cursor.Decode(&post); err != nil {
			return c.Status(http.StatusInternalServerError).SendString(err.Error())
		}
		posts = append(posts, post)
	}

	// Check if any error occurred during iteration
	if err := cursor.Err(); err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	return c.Status(http.StatusOK).JSON(allPostSerializer{Message: "Posts fetched successfully", Posts: posts})
}

func createPost(c *fiber.Ctx) error {
	PostCollection := database.Client.Database("Feed").Collection("Post")

	post := new(Post)
	if err := c.BodyParser(post); err != nil {
		return c.Status(http.StatusBadRequest).SendString(err.Error())
	}

	// Handle file upload for ImageURL
	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(Error{
			Message: "Image upload failed",
			Errors:  err.Error(),
		})
	}

	// Save the file to your server and get the URL
	// This is just an example, adjust according to your needs
	filePath := "./images/" + file.Filename
	c.SaveFile(file, filePath)
	post.ImageURL = "images/" + file.Filename

	post.SetTimestamps()
	post.Creator = creator{Name: "Jack Berry"}

	// create post in db
	result, err := PostCollection.InsertOne(context.TODO(), post)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	// Retrieve the inserted document from the database
	insertedPost := new(Post)
	err = PostCollection.FindOne(context.TODO(), bson.M{"_id": result.InsertedID}).Decode(insertedPost)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	return c.Status(http.StatusCreated).JSON(postSerializer{Message: "Post created successfully", Post: insertedPost})
}

func getPost(c *fiber.Ctx) error {
	postId := c.Params("postId")
	PostCollection := database.Client.Database("Feed").Collection("Post")

	post := new(Post)

	objectId, err := primitive.ObjectIDFromHex(postId)
	if err != nil {
		return c.Status(http.StatusBadRequest).SendString("Invalid Id")
	}

	err = PostCollection.FindOne(context.TODO(), bson.M{"_id": objectId}).Decode(post)
	if err != nil {
		// return c.Status(http.StatusBadRequest).SendString(err.Error())
		return c.Status(http.StatusBadRequest).SendString("could not find post or Invalid Id")
	}

	return c.Status(http.StatusOK).JSON(postSerializer{Message: "Post fetched successfully", Post: post})
}

func updatePost(c *fiber.Ctx) error {
	postId := c.Params("postId")
	PostCollection := database.Client.Database("Feed").Collection("Post")

	objectId, err := primitive.ObjectIDFromHex(postId)
	if err != nil {
		return c.Status(http.StatusBadRequest).SendString("Invalid Id")
	}

	oldPost := new(Post)
	err = PostCollection.FindOne(context.TODO(), bson.M{"_id": objectId}).Decode(oldPost)
	if err != nil {
		return c.Status(http.StatusBadRequest).SendString("could not find post or Invalid Id")
	}

	post := new(Post)
	if err := c.BodyParser(post); err != nil {
		return c.Status(http.StatusBadRequest).SendString(err.Error())
	}

	file, err := c.FormFile("image")
	if err != nil {
		if c.FormValue("image") == "" {
			return c.Status(http.StatusUnprocessableEntity).JSON(Error{
				Message: "No file picked.",
				Errors:  err.Error(),
			})
		}
		post.ImageURL = c.FormValue("image")
	} else {
		filePath := "./images/" + file.Filename
		c.SaveFile(file, filePath)
		post.ImageURL = "images/" + file.Filename
	}

	if post.ImageURL != oldPost.ImageURL {
		clearImage(oldPost.ImageURL)
	}

	post.SetTimestamps()

	update := bson.M{
		"$set": bson.M{
			"title":     post.Title,
			"content":   post.Content,
			"imageUrl":  post.ImageURL,
			"updatedAt": post.UpdatedAt,
			// "creator":   post.Creator,
		},
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	err = PostCollection.FindOneAndUpdate(context.TODO(), bson.M{"_id": objectId}, update, opts).Decode(post)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	slog.Info(fmt.Sprintf("post with id %s updated successfully", postId))

	return c.Status(http.StatusOK).JSON(postSerializer{Message: "Post updated successfully", Post: post})
}

func deletePost(c *fiber.Ctx) error {
	postId := c.Params("postId")
	PostCollection := database.Client.Database("Feed").Collection("Post")

	objectId, err := primitive.ObjectIDFromHex(postId)
	if err != nil {
		return c.Status(http.StatusBadRequest).SendString("Invalid Id")
	}

	deletedPost := new(Post)
	opts := options.FindOneAndDelete().SetProjection(bson.M{"_id": 1, "imageUrl": 1})
	err = PostCollection.FindOneAndDelete(context.TODO(), bson.M{"_id": objectId}, opts).Decode(deletedPost)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(http.StatusNotFound).SendString("Post not found")
		}
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	slog.Info(fmt.Sprintf("post with id %s deleted successfully", postId))

	clearImage(deletedPost.ImageURL)

	return c.Status(http.StatusOK).JSON(postSerializer{Message: "Post deleted successfully", Post: deletedPost})
}
