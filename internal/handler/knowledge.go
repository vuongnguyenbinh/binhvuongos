package handler

import (
	"binhvuongos/web/templates/pages"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) Knowledge(c *fiber.Ctx) error {
	return render(c, pages.KnowledgePage())
}

func (h *Handler) CreateKnowledgeItem(c *fiber.Ctx) error {
	// TODO: Phase 6 — parse form and create knowledge item
	return c.Redirect("/knowledge")
}
