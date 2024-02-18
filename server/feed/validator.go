package feed

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var Validator = validator.New()

func validateCreateAndUpdatePost(c *fiber.Ctx) error {
	post := new(createPostInput)

	// Check the Content-Type header
	if c.Get("Content-Type") == "application/json" {
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
			if c.FormValue("image") == "" {
				return c.Status(http.StatusBadRequest).JSON(Error{
					Message: "Image upload failed",
					Errors:  err.Error(),
				})
			}
		}

		post.Image = imageField{
			File: file,
			URL:  c.FormValue("image"),
		}
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
