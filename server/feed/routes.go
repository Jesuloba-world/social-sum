package feed

import "github.com/gofiber/fiber/v2"

func Router(app *fiber.App) {
	api := app.Group("/feed")
	api.Get("/posts", getPosts)
	api.Post("/post", createPost)
}
