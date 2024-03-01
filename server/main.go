package main

import (
	"log"
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
	"github.com/joho/godotenv"

	"github.com/Jesuloba-world/social-sum/server/auth"
	"github.com/Jesuloba-world/social-sum/server/database"
	_ "github.com/Jesuloba-world/social-sum/server/docs"
	"github.com/Jesuloba-world/social-sum/server/feed"
)

//	@title						Social sum API
//	@version					1.0
//	@description				This is the documentation for social sum api
//	@host						localhost:8000
//	@BasePath					/
//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
//	@security					[{"BearerAuth":[]}]

func main() {
	slog.Info("Application started")

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

	app.Get("/swagger/*", swagger.HandlerDefault)

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173",
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH",
		AllowHeaders:     "Content-Type, Authorization",
		AllowCredentials: true,
	}))

	auth.Router(app)
	feed.Router(app)

	app.Listen(":8000")
}
