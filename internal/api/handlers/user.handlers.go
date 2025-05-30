package handlers

import (
	"bytes"
	"fmt"
	"net/http"
	"os"

	"github.com/clerk/clerk-sdk-go/v2"
	clerkhttp "github.com/clerk/clerk-sdk-go/v2/http"
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp/fasthttpadaptor"
)

// responseRecorder is a minimal implementation of http.ResponseWriter
// that captures status codes
type responseRecorder struct {
	headers    http.Header
	statusCode int
	body       bytes.Buffer
}

// newResponseRecorder creates a new response recorder
func newResponseRecorder() *responseRecorder {
	return &responseRecorder{
		headers:    make(http.Header),
		statusCode: http.StatusOK,
	}
}

// Header returns the header map for setting HTTP headers
func (r *responseRecorder) Header() http.Header {
	return r.headers
}

// Write captures response body data
func (r *responseRecorder) Write(bytes []byte) (int, error) {
	return r.body.Write(bytes)
}

// WriteHeader captures the status code
func (r *responseRecorder) WriteHeader(statusCode int) {
	r.statusCode = statusCode
}

// WIP: modify this. Add more logic to
// THIS IS A MIDDLEWARE
func Authenticate() fiber.Handler {
	return func(c *fiber.Ctx) error {

		var CLERK_SECRET_KEY = os.Getenv("CLERK_SECRET_KEY")
		clerk.SetKey(CLERK_SECRET_KEY)

		httpReq := &http.Request{}
		ctx := c.Context()

		if err := fasthttpadaptor.ConvertRequest(ctx, httpReq, false); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"status":  http.StatusInternalServerError,
				"message": err.Error(),
			})
		}

		authHeader := c.Get("Authorization")
		if authHeader != "" {
			httpReq.Header.Set("Authorization", authHeader)
		}

		recorder := newResponseRecorder()

		httpHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := clerk.SessionClaimsFromContext(r.Context())
			fmt.Println("Subject", claims.Subject)

			if !ok {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			if claims != nil {
				c.Locals("userId", claims.Subject)
			}
		})

		clerkHandler := clerkhttp.WithHeaderAuthorization()(httpHandler)
		clerkHandler.ServeHTTP(recorder, httpReq)

		if recorder.statusCode != http.StatusOK {
			return c.Status(recorder.statusCode).JSON(fiber.Map{
				"status":  recorder.statusCode,
				"message": "Invalid session",
			})
		}

		return c.Next()
	}
}
