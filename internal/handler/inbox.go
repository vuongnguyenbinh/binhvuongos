package handler

import (
	"binhvuongos/web/templates/pages"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) Inbox(c *fiber.Ctx) error {
	return render(c, pages.InboxPage())
}

func (h *Handler) CreateInboxItem(c *fiber.Ctx) error {
	// TODO: Phase 4 — parse form and create inbox item
	return c.Redirect("/inbox")
}

func (h *Handler) TriageInbox(c *fiber.Ctx) error {
	// TODO: Phase 4 — triage inbox item
	return c.Redirect("/inbox/" + c.Params("id"))
}
