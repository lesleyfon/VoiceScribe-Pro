package main

import (
	"voicescribe-pro/internal/api/handlers"
	"voicescribe-pro/internal/api/routes"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
)

func main() {
	app := fiber.New()

	var AllowOrigins = []string{"https://www.hello.com"}

	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOrigins:     AllowOrigins, // URL for FE
		AllowHeaders:     []string{"Origin, Content-Type, Accept"},
		AllowMethods:     []string{fiber.MethodGet, fiber.MethodPost}, // REST Methods allowed
		MaxAge:           3600,                                        // How long a preflight request should be cached for

	}))
	app.Get("/unprotected", func(c fiber.Ctx) error {
		return c.JSON(map[string]string{
			"message": "Welcome to the unprotected route",
		})
	})
	app.Use(handlers.Authenticate)
	app.Get("/welcome", func(c fiber.Ctx) error {
		return c.JSON(map[string]string{
			"message": "Welcome to the home route",
		})
	})
	routes.AuthRoutes(app)

	app.Listen(":3000")
}
