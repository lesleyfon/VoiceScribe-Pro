package main

import (
	"voicescribe-pro/internal/api/handlers"
	"voicescribe-pro/internal/api/routes"

	"github.com/gofiber/fiber/v3"
)

func main() {
	app := fiber.New()

	app.Get("/unprotected", func(c fiber.Ctx) error {
		return c.JSON(map[string]string{
			"message": "Welcomme to the unprotected route",
		})
	})
	app.Use(handlers.Authenticate)
	app.Get("/welcome", func(c fiber.Ctx) error {
		return c.JSON(map[string]string{
			"message": "Welcomme to the home route",
		})
	})
	routes.AuthRoutes(app)

	app.Listen(":3000")
}
