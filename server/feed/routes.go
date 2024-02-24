package feed

import (
	"github.com/gofiber/fiber/v2"

	"github.com/Jesuloba-world/social-sum/server/middleware"
)

func Router(app *fiber.App) {
	api := app.Group("/feed")
	api.Get("/posts", middleware.IsAuth, getPosts)
	api.Post("/post", middleware.IsAuth, validateCreateAndUpdatePost, createPost)
	api.Get("/post/:postId", middleware.IsAuth, getPost)
	api.Put("/post/:postId", middleware.IsAuth, validateCreateAndUpdatePost, updatePost)
	api.Delete("/post/:postId", middleware.IsAuth, deletePost)
}
