package auth

import (
	"SoAt/internals/dto"
	"SoAt/models/users"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
)

// AuthMiddleware protects routes that require authentication
func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Authorization header required",
			})
		}

		// Check if the header starts with "Bearer "
		if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid authorization format",
			})
		}

		// Extract and validate token
		tokenString := authHeader[7:]
		userID, err := ExtractUserID(tokenString)
		if err != nil {
			// Only log in debug mode
			if os.Getenv("DEBUG") == "true" {
				fmt.Printf("Token validation error: %v\n", err)
			}
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired token",
			})
		}

		// Get user from database to verify they exist
		user := users.New()
		user.User = &dto.User{ID: userID}
		ctx := c.UserContext()
		if err := user.Get(ctx); err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "User not found",
			})
		}

		// Set user ID in context for later use
		// Store as string to avoid type assertion issues
		c.Locals("userID", userID.String())
		return c.Next()
	}
}
