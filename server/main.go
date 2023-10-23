package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"github.com/Jesuloba-world/social-sum/server/feed"
)

func main() {
	app := fiber.New(fiber.Config{
		Immutable: true,
		// EnablePrintRoutes: true,
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH",
		AllowHeaders: "Content-Type, Authorization",
	}))

	feed.Router(app)

	app.Listen(":8000")
}
