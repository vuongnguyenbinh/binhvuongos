package middleware

import (
	"crypto/subtle"

	"github.com/gofiber/fiber/v2"
)

// APIKeyAuth validates X-API-Key using a constant-time comparison to defeat timing side channels.
func APIKeyAuth(apiKey string) fiber.Handler {
	expected := []byte(apiKey)
	return func(c *fiber.Ctx) error {
		got := []byte(c.Get("X-API-Key"))
		if len(expected) == 0 || subtle.ConstantTimeCompare(got, expected) != 1 {
			return c.Status(401).JSON(fiber.Map{
				"success": false,
				"error":   "Invalid or missing API key",
			})
		}
		return c.Next()
	}
}
