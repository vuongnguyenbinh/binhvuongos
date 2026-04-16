package handler

import (
	"binhvuongos/web/templates/pages"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) Bookmarks(c *fiber.Ctx) error {
	return render(c, pages.BookmarksPage())
}

func (h *Handler) CreateBookmark(c *fiber.Ctx) error {
	// TODO: Phase 6 — parse form and create bookmark
	return c.Redirect("/bookmarks")
}

func (h *Handler) DeleteBookmark(c *fiber.Ctx) error {
	// TODO: Phase 6 — soft delete bookmark
	return c.Redirect("/bookmarks")
}
