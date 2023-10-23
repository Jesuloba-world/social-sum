package feed

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type createPostSerializer struct {
	Message string `json:"message"`
	Post    *Post  `json:"post"`
}

func getPosts(c *fiber.Ctx) error {
	return c.Status(http.StatusOK).JSON(Post{Title: "First Post", Content: "This is the first post!"})
}

func createPost(c *fiber.Ctx) error {
	post := new(Post)
	if err := c.BodyParser(post); err != nil {
		return c.Status(http.StatusBadRequest).SendString(err.Error())
	}
	// create post in db
	return c.Status(http.StatusCreated).JSON(createPostSerializer{Message: "Post created sucessfully!", Post: post})
}
