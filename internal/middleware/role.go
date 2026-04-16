package middleware

import (
	"binhvuongos/internal/db/generated"

	"github.com/gofiber/fiber/v2"
)

// RequireRole checks if the authenticated user has one of the allowed roles
func RequireRole(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, ok := c.Locals("user").(generated.User)
		if !ok {
			return c.Redirect("/login")
		}
		for _, r := range roles {
			if user.Role == r {
				return c.Next()
			}
		}
		return c.Status(403).SendString("Forbidden")
	}
}
