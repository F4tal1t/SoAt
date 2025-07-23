package routes

import (
	"SoAt/controllers/users"
	"SoAt/internals/auth"
	"github.com/gofiber/fiber/v2"
)

func Users(r fiber.Router) {
	userRoutes := r.Group("/users")
	
	// Public routes - only get specific user by ID
	userRoutes.Get("/:id", users.Get)      // Get user by ID
	
	// Protected routes - require authentication
	userRoutes.Post("/", auth.AuthMiddleware(), users.Add)        // Create user
	userRoutes.Delete("/:id", auth.AuthMiddleware(), users.Delete) // Delete user by ID
	
	// Admin-only routes
	userRoutes.Get("/", auth.AdminMiddleware(), users.GetAll)      // Get all users (admin only)
}
