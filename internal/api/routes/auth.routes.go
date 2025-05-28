package routes

import "github.com/gofiber/fiber/v3"

func AuthRoutes(app *fiber.App) {
	app.Post("/login", func(c fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})
}
