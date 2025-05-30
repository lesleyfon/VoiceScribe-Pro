package handlers

import (
	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/gofiber/fiber/v2"
)

// GetUserID extracts the user ID from Fiber context
func GetUserID(c *fiber.Ctx) (string, bool) {
	userID, ok := c.Locals("userId").(string)
	return userID, ok
}

// GetUserClaims extracts the user claims from Fiber context
func GetUserClaims(c *fiber.Ctx) (*clerk.SessionClaims, bool) {
	claims, ok := c.Locals("userClaims").(*clerk.SessionClaims)
	return claims, ok
}
