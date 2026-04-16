package handler

import (
	"binhvuongos/web/templates/pages"
	"github.com/gofiber/fiber/v2"
)

func Inbox(c *fiber.Ctx) error {
	return render(c, pages.InboxPage())
}
