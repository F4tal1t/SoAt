package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

var app *fiber.App

func New() *fiber.App {
	return app
}

func Setup() {
	app = fiber.New(fiber.Config{
		ErrorHandler: nil,
		BodyLimit:    4 * 1024 * 1024, // 4MB
	})

	defer app.Use(NotFoundHandler)
	defer app.Use(recover.New())

	Middleware(app)
	AddRoutes(app)
}
