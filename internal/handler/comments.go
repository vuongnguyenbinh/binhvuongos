package handler

import (
	"binhvuongos/internal/middleware"
	"binhvuongos/web/templates/components"

	"github.com/gofiber/fiber/v2"
)

// CreateComment handles POST /comments
func (h *Handler) CreateComment(c *fiber.Ctx) error {
	user := GetUser(c)
	entityType := c.FormValue("entity_type")
	entityID := c.FormValue("entity_id")
	body := c.FormValue("body")

	if entityType == "" || entityID == "" || body == "" {
		return c.Status(400).SendString("Thiếu thông tin")
	}

	_ = h.queries.CreateComment(c.Context(), entityType, middleware.StringToUUID(entityID), user.ID, body)

	// Return updated comment list via HTMX
	return h.renderComments(c, entityType, entityID)
}

// DeleteComment handles POST /comments/:id/delete
func (h *Handler) DeleteComment(c *fiber.Ctx) error {
	id := c.Params("id")
	entityType := c.FormValue("entity_type")
	entityID := c.FormValue("entity_id")

	_ = h.queries.DeleteComment(c.Context(), middleware.StringToUUID(id))

	if entityType != "" && entityID != "" {
		return h.renderComments(c, entityType, entityID)
	}
	return c.SendString("")
}

// LoadComments handles GET /comments?entity_type=X&entity_id=Y (HTMX partial)
func (h *Handler) LoadComments(c *fiber.Ctx) error {
	entityType := c.Query("entity_type")
	entityID := c.Query("entity_id")
	if entityType == "" || entityID == "" {
		return c.SendString("")
	}
	return h.renderComments(c, entityType, entityID)
}

func (h *Handler) renderComments(c *fiber.Ctx, entityType, entityID string) error {
	user := GetUser(c)
	comments, _ := h.queries.ListCommentsByEntity(c.Context(), entityType, middleware.StringToUUID(entityID))

	var items []components.CommentData
	for _, cm := range comments {
		items = append(items, components.CommentData{
			ID:       middleware.UUIDToString(cm.ID),
			Body:     cm.Body,
			UserName: cm.FullName,
			UserInit: string([]rune(cm.FullName)[:1]),
			TimeAgo:  timeAgo(cm.CreatedAt),
			IsOwner:  cm.UserID == user.ID,
		})
	}

	return render(c, components.CommentList(entityType, entityID, items))
}
