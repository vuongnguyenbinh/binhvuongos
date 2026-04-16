package handler

import (
	"binhvuongos/web/templates/pages"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) Content(c *fiber.Ctx) error {
	return render(c, pages.ContentPage())
}

func (h *Handler) CreateContent(c *fiber.Ctx) error {
	// TODO: Phase 5 — parse form and create content
	return c.Redirect("/content")
}

func (h *Handler) ReviewContentForm(c *fiber.Ctx) error {
	// TODO: Phase 5 — review content action
	return c.Redirect("/content/" + c.Params("id"))
}
