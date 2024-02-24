package feed

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/Jesuloba-world/social-sum/server/auth"
	"github.com/Jesuloba-world/social-sum/server/database"
)

type postSerializer struct {
	Message string  `json:"message"`
	Post    *Post   `json:"post"`
	Creator creator `json:"creator"`
}

type creator struct {
	ID   primitive.ObjectID `json:"_id"`
	Name string             `json:"name"`
}

type allPostSerializer struct {
	Message    string `json:"message"`
	Posts      []Post `json:"posts"`
	TotalItems int64  `json:"totalItems"`
}

func getPosts(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "2"))

	PostCollection := database.Client.Database("Feed").Collection("Post")

	skip := (page - 1) * limit

	opts := options.Find().SetLimit(int64(limit)).SetSkip(int64(skip))

	// Find all documents in the collection
	cursor, err := PostCollection.Find(context.TODO(), bson.M{}, opts)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	defer cursor.Close(context.TODO())

	var posts []Post

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

	// Get the total number of documents in the collection
	total, err := PostCollection.CountDocuments(context.TODO(), bson.M{})
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	return c.Status(http.StatusOK).JSON(allPostSerializer{Message: "Posts fetched successfully", Posts: posts, TotalItems: total})
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

	// add user_id as post creator
	userId, err := getUserIdFromLocals(c)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}
	post.Creator = userId

	// get user object
	UserCollection := database.Client.Database("Auth").Collection("User")
	user := new(auth.User)
	err = UserCollection.FindOne(context.TODO(), bson.M{"_id": userId}).Decode(user)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

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

	// append new post
	user.Posts = append(user.Posts, insertedPost.ID)

	user.SetTimestamps()

	// update database
	update := bson.M{
		"$set": bson.M{
			"posts":     user.Posts,
			"updatedAt": user.UpdatedAt,
		},
	}

	_, err = UserCollection.UpdateOne(context.TODO(), bson.M{"_id": user.ID}, update)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	return c.Status(http.StatusCreated).JSON(postSerializer{Message: "Post created successfully", Post: insertedPost, Creator: creator{ID: user.ID, Name: user.Name}})
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
	UserCollection := database.Client.Database("Auth").Collection("User")

	// get userId
	userId, err := getUserIdFromLocals(c)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	// get user object
	user := new(auth.User)
	err = UserCollection.FindOne(context.TODO(), bson.M{"_id": userId}).Decode(user)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	objectId, err := primitive.ObjectIDFromHex(postId)
	if err != nil {
		return c.Status(http.StatusBadRequest).SendString("Invalid Id")
	}

	deletedPost := new(Post)
	opts := options.FindOne().SetProjection(bson.M{"_id": 1, "imageUrl": 1, "creator": 1})
	err = PostCollection.FindOne(context.TODO(), bson.M{"_id": objectId}, opts).Decode(deletedPost)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(http.StatusNotFound).SendString("Post not found")
		}
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	if deletedPost.Creator != user.ID {
		return c.Status(http.StatusUnauthorized).SendString("You are not authorized to delete this post")
	}

	result, err := PostCollection.DeleteOne(context.TODO(), bson.M{"_id": deletedPost.ID})
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	if result.DeletedCount <= 0 {
		return c.Status(http.StatusInternalServerError).SendString("No document deleted")
	}

	// find the deleted post
	postIndex := -1
	for i, postId := range user.Posts {
		if postId == deletedPost.ID {
			postIndex = i
			break
		}
	}

	// if the post is in the user's post array, remove it
	if postIndex != -1 {
		user.Posts = append(user.Posts[:postIndex], user.Posts[postIndex+1:]...)

		user.SetTimestamps()

		// update database
		update := bson.M{
			"$set": bson.M{
				"posts":     user.Posts,
				"updatedAt": user.UpdatedAt,
			},
		}

		_, err = UserCollection.UpdateOne(context.TODO(), bson.M{"_id": user.ID}, update)
		if err != nil {
			return c.Status(http.StatusInternalServerError).SendString(err.Error())
		}
	}

	slog.Info(fmt.Sprintf("post with id %s deleted successfully", postId))

	clearImage(deletedPost.ImageURL)

	return c.Status(http.StatusOK).SendString("Post deleted successfully")
}
