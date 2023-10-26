package feed

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

)

var Validator = validator.New()

func validateCreatePost(c *fiber.Ctx) error {
	post := new(createPostInput)

	// Check the Content-Type header
	if c.Get("Content-Type") == "application/json" {
		println("Json")
		// Parse JSON body
		if err := c.BodyParser(post); err != nil {
			return c.Status(http.StatusBadRequest).JSON(Error{
				Message: "Validation failed, entered data is incorrect",
				Errors:  err.Error(),
			})
		}
	} else {
		// Parse form data
		post.Title = c.FormValue("title")
		post.Content = c.FormValue("content")

		// Handle file check for image
		file, err := c.FormFile("image")
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(Error{
				Message: "Image upload failed",
				Errors:  err.Error(),
			})
		}

		post.Image = file
	}

	validationErr := Validator.Struct(post)
	if validationErr != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(Error{
			Message: "Validation failed, entered data is incorrect",
			Errors:  validationErr.Error(),
		})
	}

	return c.Next()
}

func Router(app *fiber.App) {
	api := app.Group("/feed")
	api.Get("/posts", getPosts)
	api.Post("/post", validateCreatePost, createPost)
	api.Get("/post/:postId", getPost)
}
