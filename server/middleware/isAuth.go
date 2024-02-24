package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func IsAuth(c *fiber.Ctx) error {
	var secret_key = os.Getenv("SECRET_KEY")
	authHeader := c.Get("Authorization")

	var tokenString string

	if authHeader != "" {
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 || parts[0] == "Bearer" {
			tokenString = parts[1]
		}
	}

	if tokenString == "" {
		tokenString = c.Cookies("jwt")
	}

	if tokenString == "" {
		return c.Status(http.StatusUnauthorized).SendString("Token not found in Header or Cookie")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret_key), nil
	})

	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(Error{Message: "An error occured", Error: err.Error()})
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		c.Locals("user_id", claims["user_id"])
		return c.Next()
	}

	return c.Status(http.StatusUnauthorized).SendString("Invalid token")
}
