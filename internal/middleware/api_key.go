package middleware

import (
	"github.com/gofiber/fiber/v2"
)

// APIKeyAuth validates requests using X-API-Key header
func APIKeyAuth(apiKey string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		key := c.Get("X-API-Key")
		if key == "" || key != apiKey {
			return c.Status(401).JSON(fiber.Map{
				"success": false,
				"error":   "Invalid or missing API key",
			})
		}
		return c.Next()
	}
}
