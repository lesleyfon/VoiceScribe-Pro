package routes

import (
	"voicescribe-pro/internal/api/handlers"

	"github.com/gofiber/fiber/v2"
)

func ProcessAudioRoutes(app *fiber.App) {
	app.Post("/audio/process-full-audio", handlers.ProcessFullAudioHandler)
}
