package handler

import (
	"context"
	"log"

	"binhvuongos/internal/config"
	"binhvuongos/internal/db/generated"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgtype"
)

// Handler holds dependencies for all route handlers
type Handler struct {
	queries     *generated.Queries
	config      *config.Config
	ownerUserID pgtype.UUID
}

// NewHandler creates a new Handler with DB queries and config.
// Resolves the owner user ID eagerly so webhook ingest always has a valid submitted_by.
// Fail-fast: panic if OWNER_EMAIL does not match any users row — misconfigured deploy.
func NewHandler(queries *generated.Queries, cfg *config.Config) *Handler {
	h := &Handler{queries: queries, config: cfg}
	if cfg.OwnerEmail != "" {
		user, err := queries.GetUserByEmail(context.Background(), cfg.OwnerEmail)
		if err != nil {
			log.Fatalf("handler init: OWNER_EMAIL=%q not found in users table: %v", cfg.OwnerEmail, err)
		}
		h.ownerUserID = user.ID
		log.Printf("Owner user resolved: email=%s", cfg.OwnerEmail)
	}
	return h
}

// GetUser extracts the authenticated user from fiber context
func GetUser(c *fiber.Ctx) generated.User {
	user, _ := c.Locals("user").(generated.User)
	return user
}
