package routes

import (
	"SoAt/controllers/posts"
	"SoAt/internals/auth"
	"github.com/gofiber/fiber/v2"
)

func Posts(r fiber.Router) {
	postsRoutes := r.Group("/users/:id/posts")

	// Public routes
	postsRoutes.Get("/", posts.Get)
	
	// Protected routes
	postsRoutes.Post("/", auth.AuthMiddleware(), posts.Add)
	postsRoutes.Delete("/:post_id", auth.AuthMiddleware(), posts.Delete)
}
