package server

import (
	"SoAt/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func Middleware(app *fiber.App) {
	app.Use(logger.New())
}

func AddRoutes(app *fiber.App) {
	baseRouter := app.Group("/SoAt")
	routes.Users(baseRouter)
	routes.Friendships(baseRouter)
	routes.Posts(baseRouter)
}
