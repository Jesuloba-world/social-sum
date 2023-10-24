package feed

import (
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

type createPostSerializer struct {
	Message string `json:"message"`
	Post    *Post  `json:"post"`
}

func getPosts(c *fiber.Ctx) error {
	return c.Status(http.StatusOK).JSON([]Post{{
		Id:       "1",
		Title:    "First Post",
		Content:  "This is the first post!",
		ImageUrl: "images/cook.jpg",
		Creator: creator{
			Name: "John Needle",
		},
		CreatedAt: time.Now(),
	}})
}

func createPost(c *fiber.Ctx) error {
	post := new(Post)
	if err := c.BodyParser(post); err != nil {
		return c.Status(http.StatusBadRequest).SendString(err.Error())
	}

	post.CreatedAt = time.Now()
	post.Creator = creator{Name: "Jack Berry"}
	// create post in db
	return c.Status(http.StatusCreated).JSON(createPostSerializer{Message: "Post created sucessfully!", Post: post})
}
