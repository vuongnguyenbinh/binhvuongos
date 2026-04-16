package handler

import (
	"binhvuongos/web/templates/pages"
	"github.com/gofiber/fiber/v2"
)

func Bookmarks(c *fiber.Ctx) error {
	return render(c, pages.BookmarksPage())
}
