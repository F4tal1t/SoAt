package routes

import (
	"SoAt/controllers/users"
	"github.com/gofiber/fiber/v2"
)

func Users(r fiber.Router) {
	userRoutes := r.Group("/users")
	userRoutes.Post("/", users.Add)        // Create user
	userRoutes.Get("/", users.GetAll)      // Get all users
	userRoutes.Get("/:id", users.Get)      // Get user by ID
	userRoutes.Delete("/:id", users.Delete) // Delete user by ID
}
