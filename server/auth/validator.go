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

func validateSignup(c *fiber.Ctx) error {
	input := new(signupInput)

	if err := c.BodyParser(input); err != nil {
		return c.Status(http.StatusBadRequest).JSON(Error{
			Message: "Validation failed, entered data is incorrect",
			Error:   err.Error(),
		})
	}

	validationErr := Validator.Struct(input)

	if validationErr != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(Error{
			Message: "Validation failed",
			Error:   validationErr.Error(),
		})
	}

	// check
	userCollection := database.Client.Database("Auth").Collection("User")
	filter := bson.M{"email": input.Email}
	result := new(User)
	err := userCollection.FindOne(context.TODO(), filter).Decode(result)

	if err == nil && result.Email != "" {
		return c.Status(http.StatusBadRequest).SendString("The User already exists")
	}

	return c.Next()
}

func validateLogin(c *fiber.Ctx) error {
	input := new(loginInput)

	if err := c.BodyParser(input); err != nil {
		return c.Status(http.StatusBadRequest).JSON(Error{
			Message: "An error occured",
			Error:   err.Error(),
		})
	}

	validationErr := Validator.Struct(input)

	if validationErr != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(Error{
			Message: "Validation failed",
			Error:   validationErr.Error(),
		})
	}

	return c.Next()
}
