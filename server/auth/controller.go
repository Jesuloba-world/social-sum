package auth

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Jesuloba-world/social-sum/server/database"
)

type userSerializer struct {
	Message string `json:"message"`
	UserID  string `json:"userid"`
}

func signup(c *fiber.Ctx) error {
	userCollection := database.Client.Database("Auth").Collection("User")

	input := new(signupInput)

	if err := c.BodyParser(input); err != nil {
		return c.Status(http.StatusBadRequest).SendString(err.Error())
	}

	hashedPassword, err := hashPassword(input.Password)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	user := User{
		Email:    input.Email,
		Name:     input.Name,
		Password: hashedPassword,
		Status:   "I am new!",
	}

	user.SetTimestamps()

	// create user in db
	result, err := userCollection.InsertOne(context.TODO(), user)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	userID := result.InsertedID.(primitive.ObjectID).Hex()

	slog.Info("User created!")

	return c.Status(http.StatusOK).JSON(userSerializer{Message: "User created successfully", UserID: userID})
}
