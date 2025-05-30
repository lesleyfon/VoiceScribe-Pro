package routes

import (
	"voicescribe-pro/internal/api/handlers"

	"github.com/gofiber/fiber/v2"
)

func UserRoutes(app *fiber.App) {
	app.Get("/user-info", handlers.GetUserClaimsHandler)
}
