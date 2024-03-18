package main

import (
	"log"
	"log/slog"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
	"github.com/joho/godotenv"
	"github.com/valyala/fasthttp/fasthttpadaptor"

	"github.com/Jesuloba-world/social-sum/server/auth"
	"github.com/Jesuloba-world/social-sum/server/database"
	_ "github.com/Jesuloba-world/social-sum/server/docs"
	"github.com/Jesuloba-world/social-sum/server/feed"
	"github.com/Jesuloba-world/social-sum/server/graph"
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

func wrapHandler(f func(http.ResponseWriter, *http.Request)) func(ctx *fiber.Ctx) {
	return func(ctx *fiber.Ctx) {
		fasthttpadaptor.NewFastHTTPHandler(http.HandlerFunc(f))(ctx.Context())
	}
}

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

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		DB: database.Client,
	}}))

	// Serve GraphQL API
	app.Use("/graphql", func(c *fiber.Ctx) error {
		wrapHandler(srv.ServeHTTP)(c)
		return nil
	})

	// Serve GraphQL Playground
	app.Get("/playground", func(c *fiber.Ctx) error {
		wrapHandler(playground.Handler("GraphQL", "/graphql"))(c)
		return nil
	})

	app.Static("/images", "./images")

	app.Get("/swagger/*", swagger.HandlerDefault)

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173,https://altair-gql.sirmuel.design",
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH",
		AllowHeaders:     "Content-Type, Authorization",
		AllowCredentials: true,
	}))

	auth.Router(app)
	feed.Router(app)
	app.Get("/ws", websocket.New(feed.BroadcastHandler))

	app.Listen(":8000")
}
