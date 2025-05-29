package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/user"
	"github.com/gofiber/fiber/v3"
)

// WIP: modify this. Add more logic to
// THIS IS A MIDDLEWARE
func Authenticate(c fiber.Ctx) error {

	ctx := context.Background()

	claims, ok := clerk.SessionClaimsFromContext(c.Context())

	if !ok {
		var code = http.StatusUnauthorized
		return c.JSON(map[string]string{
			"statusCode": fmt.Sprintf("%d", code),
			"access":     "unauthorized",
		})
	}

	usr, err := user.Get(ctx, claims.Subject)

	if err != nil {
		var code = http.StatusUnauthorized
		return c.JSON(map[string]string{
			"statusCode": fmt.Sprintf("%d", code),
			"error":      err.Error(),
		})
	}

	fmt.Printf(`{"user_id": "%s", "user_banned": "%t"}`, usr.ID, usr.Banned)
	return c.Next()
}
