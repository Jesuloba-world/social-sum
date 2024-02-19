package auth

import "github.com/gofiber/fiber/v2"

func Router(app *fiber.App) {
	api := app.Group("/auth")
	api.Post("/signup", validateSignup, signup)
	api.Post("/login", validateLogin, login)
}
