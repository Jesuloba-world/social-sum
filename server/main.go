package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"

	"github.com/Jesuloba-world/social-sum/server/database"
	"github.com/Jesuloba-world/social-sum/server/feed"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	app := fiber.New(fiber.Config{
		Immutable: true,
		// EnablePrintRoutes: true,
	})

	disconnect := database.Connect()

	defer disconnect()

	app.Static("/images", "./images")

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH",
		AllowHeaders: "Content-Type, Authorization",
	}))

	feed.Router(app)

	app.Listen(":8000")
}
