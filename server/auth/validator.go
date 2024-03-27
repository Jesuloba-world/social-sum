package auth

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/Jesuloba-world/social-sum/server/database"
)

var Validator = validator.New()

func validateSignup(c *fiber.Ctx) error {
	input := new(SignupInput)

	if err := c.BodyParser(input); err != nil {
		return c.Status(http.StatusBadRequest).JSON(Error{
			Message: "Validation failed, entered data is incorrect",
			Error:   err.Error(),
		})
	}

	err := ValidateSignupInput(*input)
	if err != nil {
		c.Status(http.StatusUnprocessableEntity).JSON(Error{
			Message: "Validation failed",
			Error:   err.Error(),
		})
	}

	return c.Next()
}

func ValidateSignupInput(input SignupInput) error {
	validationErr := Validator.Struct(input)

	if validationErr != nil {
		return fmt.Errorf("%s", validationErr.Error())
	}

	// check
	userCollection := database.Client.Database("Auth").Collection("User")
	filter := bson.M{"email": input.Email}
	result := new(User)
	err := userCollection.FindOne(context.TODO(), filter).Decode(result)

	if err == nil && result.Email != "" {
		return fmt.Errorf("validation failed: %s", "This User already exists")
	}

	return nil
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
