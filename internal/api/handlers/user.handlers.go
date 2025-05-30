package handlers

import (
	"bytes"
	"log"
	"net/http"
	"os"
	"time"

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

// AuthError represents different types of authentication errors
type AuthError struct {
	Type    string
	Message string
	Code    int
}

func (e AuthError) Error() string {
	return e.Message
}

var (
	ErrMissingToken = AuthError{
		Type:    "missing_token",
		Message: "Authorization token is required",
		Code:    fiber.StatusUnauthorized,
	}
	ErrInvalidToken = AuthError{
		Type:    "invalid_token",
		Message: "Invalid or expired token",
		Code:    fiber.StatusUnauthorized,
	}
	ErrInsufficientPermissions = AuthError{
		Type:    "insufficient_permissions",
		Message: "Insufficient permissions for this resource",
		Code:    fiber.StatusForbidden,
	}
	ErrInternalError = AuthError{
		Type:    "internal_error",
		Message: "Internal authentication error",
		Code:    fiber.StatusInternalServerError,
	}
)

// extractAuthToken extracts the auth token from various sources
// @param c *fiber.Ctx
// @return string
func extractAuthToken(c *fiber.Ctx) string {
	// Try Authorization header first
	if auth := c.Get("Authorization"); auth != "" {
		return auth
	}

	// Try query parameter as fallback
	if token := c.Query("token"); token != "" {
		return "Bearer " + token
	}

	// Try cookie as another fallback
	if sessionCookie := c.Cookies("__session"); sessionCookie != "" {
		return "Bearer " + sessionCookie
	}

	return ""
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

// Authenticate is a middleware that authenticates the user
func Authenticate() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		var CLERK_SECRET_KEY = os.Getenv("CLERK_SECRET_KEY")
		clerk.SetKey(CLERK_SECRET_KEY)

		httpReq := &http.Request{}
		ctx := c.Context()

		if err := fasthttpadaptor.ConvertRequest(ctx, httpReq, false); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(ErrInternalError)
		}

		authHeader := extractAuthToken(c)
		if authHeader == "" {
			return c.Status(http.StatusUnauthorized).JSON(ErrMissingToken)
		}

		// Set the Authorization header in the request
		httpReq.Header.Set("Authorization", authHeader)

		// Create a new response recorder
		recorder := newResponseRecorder()
		httpHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := clerk.SessionClaimsFromContext(r.Context())

			if !ok {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			if claims != nil {
				c.Locals("userId", claims.Subject)
				c.Locals("userClaims", claims)

				log.Println("Authentication successful", fiber.Map{
					"user_id":  claims.Subject,
					"path":     c.Path(),
					"duration": time.Since(start),
				})
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
