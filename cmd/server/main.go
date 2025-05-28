package main

import (
	"voicescribe-pro/internal/api/routes"

	"github.com/gofiber/fiber/v3"
)

func main() {
	app := fiber.New()

	routes.HelloWorld(app)
	routes.AuthRoutes(app)

	app.Listen(":3000")
}
