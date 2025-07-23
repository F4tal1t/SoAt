package routes

import (
	"SoAt/controllers/posts"
	"github.com/gofiber/fiber/v2"
)

func Posts(r fiber.Router) {
	postsRoutes := r.Group("/users/:id/posts")

	postsRoutes.Post("/", posts.Add)
	postsRoutes.Get("/", posts.Get)
	postsRoutes.Delete("/:post_id", posts.Delete)
}
