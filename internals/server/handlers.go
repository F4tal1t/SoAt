package server

import "github.com/gofiber/fiber/v2"

// Default Error Handler
func ErrorHandler(c *fiber.Ctx, err error) error {
	msg := err.Error()
	return c.Status(fiber.StatusInternalServerError).JSON(msg)
}

// Not Found Error Handler
func NotFoundHandler(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotFound).JSON("Not Found")
}
