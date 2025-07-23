package auth

import (
	"SoAt/models/users"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// AdminMiddleware checks if the authenticated user is an admin
func AdminMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Verify authentication and admin role

		// Get the Authorization header
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

		tokenString := authHeader[7:]

		// Validate the token and extract user ID
		token, err := ValidateToken(tokenString)
		if err != nil {
			// Only log in debug mode
			if os.Getenv("DEBUG") == "true" {
				fmt.Printf("Token validation error: %v\n", err)
			}
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired token",
			})
		}

		// Extract claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token claims",
			})
		}

		// Get user ID from claims
		idStr, ok := claims["id"].(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid user ID in token",
			})
		}

		// Parse the UUID
		userID, err := uuid.Parse(idStr)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid user ID format",
			})
		}

		// Get the user from the database directly
		userModel := users.Users{ID: userID}
		ctx := c.UserContext()
		if err := userModel.Get(ctx); err != nil {
			// Only log in debug mode
			if os.Getenv("DEBUG") == "true" {
				fmt.Printf("Error getting user: %v\n", err)
			}
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "User not found",
			})
		}

		// Check if the user is an admin
		if userModel.Role != "admin" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Admin access required",
			})
		}

		// Store the user ID in context for later use
		c.Locals("userID", userID)
		return c.Next()
	}
}
