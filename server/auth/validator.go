package auth

import (
	"context"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/Jesuloba-world/social-sum/server/database"
)

var Validator = validator.New()

type signupInput struct {
	Email    string `json:"email" validate:"required,email"`
	Name     string `json:"name" validate:"required"`
	Password string `json:"password" validate:"required,min=5"`
}

func validateSignup(c *fiber.Ctx) error {
	input := new(signupInput)

	if err := c.BodyParser(input); err != nil {
		return c.Status(http.StatusBadRequest).JSON(Error{
			Message: "Validation failed, entered data is incorrect",
			Errors:  err.Error(),
		})
	}

	validationErr := Validator.Struct(input)

	if validationErr != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(Error{
			Message: "Validation failed",
			Errors:  validationErr.Error(),
		})
	}

	userCollection := database.Client.Database("Auth").Collection("User")
	filter := bson.M{"email": input.Email}
	var result bson.M
	err := userCollection.FindOne(context.TODO(), filter).Decode(result)

	if err == nil {
		return c.Status(http.StatusBadRequest).SendString("Email already exists")
	}

	return c.Next()
}
