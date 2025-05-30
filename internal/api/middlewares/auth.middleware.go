package middlewares

import (
	"bytes"
	"errors"
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

func validateToken(req *http.Request) (*clerk.SessionClaims, error) {

	recorder := newResponseRecorder()
	var claims *clerk.SessionClaims
	var validationErr error

	httpHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		extractedClaims, ok := clerk.SessionClaimsFromContext(r.Context())

		if !ok || extractedClaims == nil {
			validationErr = errors.New("failed to extract session claims")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		claims = extractedClaims
	})

	clerkHandler := clerkhttp.WithHeaderAuthorization()(httpHandler)
	clerkHandler.ServeHTTP(recorder, req)

	log.Println("Recorder status code", recorder.statusCode)
	log.Println("Recorder status code is not OK", validationErr)
	if recorder.statusCode != http.StatusOK {
		return nil, validationErr
	}
	return claims, nil
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

		authToken := extractAuthToken(c)
		if authToken == "" {
			return c.Status(http.StatusUnauthorized).JSON(ErrMissingToken)
		}

		// Set the Authorization header in the request
		httpReq.Header.Set("Authorization", authToken)
		claims, err := validateToken(httpReq)

		if err != nil {
			log.Println("Token validation failed", fiber.Map{
				"error": err,
				"path":  c.Path(),
				"ip":    c.IP(),
			})
			return c.Status(http.StatusUnauthorized).JSON(ErrInvalidToken)
		}

		if claims == nil {
			log.Println("Claims are nil")
			return c.Status(http.StatusUnauthorized).JSON(ErrMissingToken)
		}

		c.Locals("userId", claims.Subject)
		c.Locals("userClaims", claims)

		log.Println("Authentication successful", fiber.Map{
			"user_id":  claims.Subject,
			"path":     c.Path(),
			"duration": time.Since(start),
		})

		return c.Next()
	}
}
