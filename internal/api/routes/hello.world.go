package routes

import "github.com/gofiber/fiber/v3"

func HelloWorld(app *fiber.App) {
	// Test route
	app.Get("/", func(c fiber.Ctx) error {
		// Send a string response to the client
		return c.SendString("Hello, World 👋!")
	})
}
