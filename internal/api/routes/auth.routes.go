package routes

import (
	"voicescribe-pro/internal/api/handlers"

	"github.com/gofiber/fiber/v3"
)

// TODO: WE DO NOT NEED THIS SINCE WE ARE USING CLERK FOR AUTH
func AuthRoutes(app *fiber.App) {
	app.Post("/authenticate", handlers.Authenticate)
}
