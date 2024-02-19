package auth

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"

	"github.com/Jesuloba-world/social-sum/server/database"
)

const SecretKey = "4nJ6rP8V9z3x7y1O4m2N6pQ8K3L6M9N0"

type userSerializer struct {
	Message string `json:"message"`
	UserID  string `json:"userid"`
}

type loginSerializer struct {
	Token  string `json:"token"`
	UserID string `json:"userid"`
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

func login(c *fiber.Ctx) error {
	userCollection := database.Client.Database("Auth").Collection("User")

	input := new(loginInput)

	if err := c.BodyParser(input); err != nil {
		return c.Status(http.StatusBadRequest).SendString(err.Error())
	}

	user := new(User)

	err := userCollection.FindOne(context.TODO(), bson.M{"email": input.Email}).Decode(user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(http.StatusNotFound).SendString("No user found with the provided email")
		}
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	// compare password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		return c.Status(http.StatusUnauthorized).SendString("Invalid Email or password")
	}

	// create claim
	expirationTime := time.Now().Add(1 * time.Hour) // 1 hour
	claim := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID.Hex(),
		"email":   user.Email,
		"exp":     expirationTime.Unix(),
	})

	// sign the the claim
	token, err := claim.SignedString([]byte(SecretKey))

	if err != nil {
		slog.Error(fmt.Sprintf("could not login: %s", err.Error()))
		return c.Status(http.StatusInternalServerError).SendString("could not login")
	}

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  expirationTime,
		HTTPOnly: true,
		SameSite: "None",
		Secure:   true,
	}

	c.Cookie(&cookie)

	return c.Status(http.StatusOK).JSON(loginSerializer{Token: token, UserID: user.ID.Hex()})
}
