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

	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOrigins:     AllowOrigins, // URL for FE
		AllowHeaders:     AllowHeaders,
		AllowMethods:     AllowMethods, // REST Methods allowed
		MaxAge:           3600,         // How long a preflight request should be cached for

	}))
	app.Get("/unprotected", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Welcome to the unprotected route",
		})
	})
	app.Use(middlewares.Authenticate())

	app.Get("/welcome", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Welcome to the home route",
		})
	})
	routes.UserRoutes(app)

	// Add fallback routes

	app.Listen(":8000")

}
