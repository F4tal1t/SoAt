package routes

import (
	"SoAt/controllers/friendships"
	"SoAt/internals/auth"
	"github.com/gofiber/fiber/v2"
)

func Friendships(r fiber.Router) {
	friendRoutes := r.Group("/friendships")

	// All friendship routes require authentication
	friendRoutes.Post("/", auth.AuthMiddleware(), friendships.Add)                         // Create friendship
	friendRoutes.Get("/:id", auth.AuthMiddleware(), friendships.Get)                       // Get friendship by ID
	friendRoutes.Get("/user/:userId", auth.AuthMiddleware(), friendships.GetAllFriends)    // Get all friends of a user
	friendRoutes.Delete("/:id", auth.AuthMiddleware(), friendships.Delete)                 // Delete friendship
}
