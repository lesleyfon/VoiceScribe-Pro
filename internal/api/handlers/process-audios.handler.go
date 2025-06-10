package handlers

import (
	"log"
	"voicescribe-pro/internal/services"

	"github.com/gofiber/fiber/v2"
)

type AudioFile struct {
	FileName    string `json:"filename"`
	FileContent string `json:"filecontnet"`
}

func ProcessFullAudioHandler(c *fiber.Ctx) error {

	if c.Method() != fiber.MethodPost && c.Method() != fiber.MethodPut {
		return c.Status(fiber.StatusMethodNotAllowed).JSON(fiber.Map{
			"error": "Method not allowed",
		})
	}

	fileHeader, err := c.FormFile("audio-file")

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"error": "Failed to get file: " + err.Error(),
			},
		)
	}

	file, err := fileHeader.Open()

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to open file: " + err.Error(),
		})
	}
	defer file.Close()

	buffer := make([]byte, fileHeader.Size)
	_, err = file.Read(buffer)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"error": "Failed to read file: " + err.Error(),
			})
	}

	transcriptionResponseData, err := services.ProcessAudio(buffer)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Transcription failed: " + err.Error(),
		})
	}

	log.Print("Body Completed")
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": "completed",
		"data":   transcriptionResponseData,
	})
}
