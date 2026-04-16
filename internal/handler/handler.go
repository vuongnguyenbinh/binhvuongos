package handler

import (
	"binhvuongos/internal/config"
	"binhvuongos/internal/db/generated"

	"github.com/gofiber/fiber/v2"
)

// Handler holds dependencies for all route handlers
type Handler struct {
	queries *generated.Queries
	config  *config.Config
}

// NewHandler creates a new Handler with DB queries and config
func NewHandler(queries *generated.Queries, cfg *config.Config) *Handler {
	return &Handler{
		queries: queries,
		config:  cfg,
	}
}

// GetUser extracts the authenticated user from fiber context
func GetUser(c *fiber.Ctx) generated.User {
	user, _ := c.Locals("user").(generated.User)
	return user
}
