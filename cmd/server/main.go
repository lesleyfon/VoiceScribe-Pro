package main

import (
	"fmt"
	"log"
	"strings"
	"voicescribe-pro/internal/api/middlewares"
	"voicescribe-pro/internal/api/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading env variables")
	}
	app := fiber.New()

	var AllowOrigins = strings.Join([]string{"http://localhost:3000"}, ",")
	fmt.Println(AllowOrigins)
	var AllowHeaders = strings.Join([]string{"Origin, Content-Type, Accept, Authorization"}, ",")
	var AllowMethods = strings.Join([]string{fiber.MethodGet, fiber.MethodPost}, ",")

	// Middlewares
	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOrigins:     AllowOrigins, // URL for FE
		AllowHeaders:     AllowHeaders,
		AllowMethods:     AllowMethods, // REST Methods allowed
		MaxAge:           3600,         // How long a preflight request should be cached for

	}))

	app.Use(middlewares.Authenticate())

	// Routes
	routes.UserRoutes(app)

	// Add fallback routes
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  fiber.StatusNotFound,
			"message": "Route not found",
			"path":    c.Path(),
		})
	})

	app.Listen(":8000")

}
