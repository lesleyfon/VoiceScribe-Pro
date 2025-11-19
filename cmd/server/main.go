package main

import (
	"log"
	"strings"
	"time"
	"voicescribe-pro/internal/api/middlewares"
	"voicescribe-pro/internal/api/routes"
	"voicescribe-pro/internal/db"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Print("Error loading env variables")
	}

	database, err := db.InitDB()
	if err != nil {
		log.Fatal("Error while connecting to DB: ", err)
	}

	app := fiber.New(
		fiber.Config{
			BodyLimit: 200 * 1024 * 1024, // 10MB in bytes
		},
	)

	var AllowOrigins = strings.Join([]string{"http://localhost:3000"}, ",")
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

	app.Use(logger.New(logger.Config{
		Format:     "${time} ${status} ${method} ${path}\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "America/Chicago",
	}))

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": "healthy",
		})
	})

	type NoteResponse struct {
		ID        uint      `json:"ID"`
		CreatedAt time.Time `json:"CreatedAt"`
		Title     string    `json:"Title"`
	}

	app.Get("/notes", func(c *fiber.Ctx) error {
		var notes []NoteResponse
		result := database.Model(&db.NotesStruct{}).Select("id", "created_at", "title").Find(&notes)

		if result.Error != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status": "error",
				"error":  result.Error.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{"notes": notes})
	})
	app.Use(middlewares.Authenticate())

	// Websocket request
	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		// Handle incoming messages
		for {
			_, msg, err := c.ReadMessage()
			if err != nil {
				log.Println("read error:", err)
				break
			}

			log.Printf("Received: %s", msg)

			// Send JSON response
			response := map[string]string{
				"pong": "pong",
			}

			if err := c.WriteJSON(response); err != nil {
				log.Println("write error:", err)
				break
			}
		}
	}))

	// Routes
	routes.UserRoutes(app)
	routes.ProcessAudioRoutes(app)

	allRoutes := app.GetRoutes()

	for _, route := range allRoutes {
		if route.Path == "/" || route.Method == "HEAD" {
			continue
		}
		log.Printf("Method: %s, Path: %s\n", route.Method, route.Path)
	}

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
