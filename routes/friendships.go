package routes

import (
	"SoAt/controllers/friendships"
	"github.com/gofiber/fiber/v2"
)

func Friendships(r fiber.Router) {
	friendRoutes := r.Group("/friendships")

	friendRoutes.Post("/", friendships.Add)                         // Create friendship
	friendRoutes.Get("/:id", friendships.Get)                       // Get friendship by ID
	friendRoutes.Get("/user/:userId", friendships.GetAllFriends)    // Get all friends of a user
	friendRoutes.Delete("/:id", friendships.Delete)                 // Delete friendship
}
