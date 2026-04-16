package handler

import (
	"binhvuongos/web/templates/pages"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) Campaigns(c *fiber.Ctx) error {
	return render(c, pages.CampaignsPage())
}

func (h *Handler) CreateCampaign(c *fiber.Ctx) error {
	// TODO: Phase 6 — parse form and create campaign
	return c.Redirect("/campaigns")
}
