package handler

import (
	"binhvuongos/web/templates/pages"
	"github.com/gofiber/fiber/v2"
)

func Campaigns(c *fiber.Ctx) error {
	return render(c, pages.CampaignsPage())
}
