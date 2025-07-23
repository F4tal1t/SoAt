package routes

import (
	"SoAt/controllers/auth"
	"github.com/gofiber/fiber/v2"
)

func Auth(r fiber.Router) {
	authRoutes := r.Group("/auth")
	authRoutes.Post("/register", auth.Register)
	authRoutes.Post("/login", auth.Login)
}