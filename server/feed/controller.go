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
	Message string   `json:"message"`
	Post    *Post    `json:"post"`
	Creator *creator `json:"creator"`
}

type allPostSerializer struct {
	Message    string `json:"message"`
	Posts      []Post `json:"posts"`
	TotalItems int64  `json:"totalItems"`
}

// @Summary		Get all posts
// @Description	Fetches all posts with pagination
// @Tags			Feed
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Param			page	query		int					false	"Page number"
// @Param			limit	query		int					false	"Number of posts per page"
// @Success		200		{object}	allPostSerializer	"Successfully fetched posts"
// @Failure		401		{string}	string				"Unauthorized"
// @Failure		500		{string}	string				"Internal Server Error"
// @Router			/feed/posts [get]
func getPosts(c *fiber.Ctx) error {
	postCollection := database.Client.Database("Feed").Collection("Post")
	userCollection := database.Client.Database("Auth").Collection("User")

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "2"))

	skip := (page - 1) * limit

	opts := options.Find().SetLimit(int64(limit)).SetSkip(int64(skip)).SetSort(bson.D{{Key: "createdAt", Value: -1}})

	// Find all documents in the collection
	cursor, err := postCollection.Find(context.TODO(), bson.M{}, opts)
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
		user := new(auth.User)
		err = userCollection.FindOne(context.TODO(), bson.M{"_id": post.CreatorId}).Decode(user)
		if err != nil {
			return c.Status(http.StatusInternalServerError).SendString(err.Error())
		}
		post.Creator = creator{Name: user.Name}
		posts = append(posts, post)
	}

	// Check if any error occurred during iteration
	if err := cursor.Err(); err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	// Get the total number of documents in the collection
	total, err := postCollection.CountDocuments(context.TODO(), bson.M{})
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	return c.Status(http.StatusOK).JSON(allPostSerializer{Message: "Posts fetched successfully", Posts: posts, TotalItems: total})
}

// @Summary		Create a new post
// @Description	Create a new post with an image and associate it with the authenticated user
// @Tags			Feed
// @Accept			json
// @Produce		json
// @Param			image	formData	file	true	"Image file"
// @Param			title	formData	string	true	"Title of the post"
// @Param			content	formData	string	true	"Content of the post"
// @Security		BearerAuth
// @Success		201	{object}	postSerializer	"Post created successfully"
// @Failure		400	{string}	string			"Bad Request"
// @Failure		500	{string}	string			"Internal Server Error"
// @Router			/feed/post [post]
func createPost(c *fiber.Ctx) error {
	postCollection := database.Client.Database("Feed").Collection("Post")
	userCollection := database.Client.Database("Auth").Collection("User")

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
	post.CreatorId = userId

	// get user object
	user := new(auth.User)
	err = userCollection.FindOne(context.TODO(), bson.M{"_id": userId}).Decode(user)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	// create post in db
	result, err := postCollection.InsertOne(context.TODO(), post)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	// Retrieve the inserted document from the database
	insertedPost := new(Post)
	err = postCollection.FindOne(context.TODO(), bson.M{"_id": result.InsertedID}).Decode(insertedPost)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	// append new post to user
	user.Posts = append(user.Posts, insertedPost.ID)

	user.SetTimestamps()

	// update database
	update := bson.M{
		"$set": bson.M{
			"posts":     user.Posts,
			"updatedAt": user.UpdatedAt,
		},
	}

	_, err = userCollection.UpdateOne(context.TODO(), bson.M{"_id": user.ID}, update)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	insertedPost.Creator = creator{Name: user.Name}

	broadcastPost(broadcastPostType{Action: "create", Post: insertedPost})

	return c.Status(http.StatusCreated).JSON(postSerializer{Message: "Post created successfully", Post: insertedPost, Creator: &creator{Name: user.Name}})
}

// @Summary		Get a specific post
// @Description	Fetches a specific post by its ID
// @Tags			Feed
// @Accept			json
// @Produce		json
// @Param			postId	path	string	true	"Post ID"
// @Security		BearerAuth
// @Success		200	{object}	postSerializer	"Post fetched successfully"
// @Failure		400	{string}	string			"Bad Request"
// @Failure		500	{string}	string			"Internal Server Error"
// @Router			/feed/post/{postId} [get]
func getPost(c *fiber.Ctx) error {
	postCollection := database.Client.Database("Feed").Collection("Post")
	userCollection := database.Client.Database("Auth").Collection("User")

	postId := c.Params("postId")

	objectId, err := primitive.ObjectIDFromHex(postId)
	if err != nil {
		return c.Status(http.StatusBadRequest).SendString("Invalid Id")
	}

	post := new(Post)

	err = postCollection.FindOne(context.TODO(), bson.M{"_id": objectId}).Decode(post)
	if err != nil {
		// return c.Status(http.StatusBadRequest).SendString(err.Error())
		return c.Status(http.StatusBadRequest).SendString("could not find post or Invalid Id")
	}

	user := new(auth.User)
	err = userCollection.FindOne(context.TODO(), bson.M{"_id": post.CreatorId}).Decode(user)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	post.Creator = creator{Name: user.Name}

	return c.Status(http.StatusOK).JSON(postSerializer{Message: "Post fetched successfully", Post: post})
}

// @Summary		Update a specific post
// @Description	Update the details of a specific post by its ID
// @Tags			Feed
// @Accept			json
// @Produce		json
// @Param			postId	path		string	true	"Post ID"
// @Param			image	formData	file	true	"Image file"
// @Param			title	formData	string	true	"Title of the post"
// @Param			content	formData	string	true	"Content of the post"
// @Security		BearerAuth
// @Success		200	{object}	postSerializer	"Post updated successfully"
// @Failure		400	{string}	string			"Bad Request"
// @Failure		401	{string}	string			"Unauthorized"
// @Failure		500	{string}	string			"Internal Server Error"
// @Router			/feed/post/{postId} [put]
func updatePost(c *fiber.Ctx) error {
	postCollection := database.Client.Database("Feed").Collection("Post")
	userCollection := database.Client.Database("Auth").Collection("User")

	postId := c.Params("postId")

	objectId, err := primitive.ObjectIDFromHex(postId)
	if err != nil {
		return c.Status(http.StatusBadRequest).SendString("Invalid Id")
	}

	oldPost := new(Post)
	err = postCollection.FindOne(context.TODO(), bson.M{"_id": objectId}).Decode(oldPost)
	if err != nil {
		return c.Status(http.StatusBadRequest).SendString("could not find post or Invalid Id")
	}

	// get userId
	userId, err := getUserIdFromLocals(c)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	if oldPost.CreatorId != userId {
		return c.Status(http.StatusUnauthorized).SendString("Not authorized!")
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
	err = postCollection.FindOneAndUpdate(context.TODO(), bson.M{"_id": objectId}, update, opts).Decode(post)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	user := new(auth.User)
	err = userCollection.FindOne(context.TODO(), bson.M{"_id": post.CreatorId}).Decode(user)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	post.Creator = creator{Name: user.Name}

	slog.Info(fmt.Sprintf("post with id %s updated successfully", postId))

	broadcastPost(broadcastPostType{Action: "update", Post: post})

	return c.Status(http.StatusOK).JSON(postSerializer{Message: "Post updated successfully", Post: post})
}

// @Summary		Delete a specific post
// @Description	Deletes a specific post by its ID
// @Tags			Feed
// @Accept			json
// @Produce		json
// @Param			postId	path	string	true	"Post ID"
// @Security		BearerAuth
// @Success		200	{string}	string	"Post deleted successfully"
// @Failure		400	{string}	string	"Bad Request"
// @Failure		401	{string}	string	"Unauthorized"
// @Failure		500	{string}	string	"Internal Server Error"
// @Router			/feed/post/{postId} [delete]
func deletePost(c *fiber.Ctx) error {
	postId := c.Params("postId")

	postCollection := database.Client.Database("Feed").Collection("Post")
	userCollection := database.Client.Database("Auth").Collection("User")

	// get userId
	userId, err := getUserIdFromLocals(c)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	// get user object
	user := new(auth.User)
	err = userCollection.FindOne(context.TODO(), bson.M{"_id": userId}).Decode(user)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	objectId, err := primitive.ObjectIDFromHex(postId)
	if err != nil {
		return c.Status(http.StatusBadRequest).SendString("Invalid Id")
	}

	deletedPost := new(Post)
	opts := options.FindOne().SetProjection(bson.M{"_id": 1, "imageUrl": 1, "creator": 1})
	err = postCollection.FindOne(context.TODO(), bson.M{"_id": objectId}, opts).Decode(deletedPost)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(http.StatusNotFound).SendString("Post not found")
		}
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	if deletedPost.CreatorId != user.ID {
		return c.Status(http.StatusUnauthorized).SendString("You are not authorized to delete this post")
	}

	result, err := postCollection.DeleteOne(context.TODO(), bson.M{"_id": deletedPost.ID})
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	clearImage(deletedPost.ImageURL)

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

		_, err = userCollection.UpdateOne(context.TODO(), bson.M{"_id": user.ID}, update)
		if err != nil {
			return c.Status(http.StatusInternalServerError).SendString(err.Error())
		}
	}

	slog.Info(fmt.Sprintf("post with id %s deleted successfully", postId))

	broadcastPost(broadcastPostType{Action: "delete", Post: deletedPost})

	return c.Status(http.StatusOK).JSON(postSerializer{Message: "Post deleted successfully", Post: deletedPost})
}
