package handler

import (
	"binhvuongos/web/templates/pages"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) ContentDetail(c *fiber.Ctx) error {
	return render(c, pages.ContentDetailPage())
}

func (h *Handler) CompanyDetail(c *fiber.Ctx) error {
	return render(c, pages.CompanyDetailPage())
}

func (h *Handler) TaskDetail(c *fiber.Ctx) error {
	return render(c, pages.TaskDetailPage())
}

func (h *Handler) WorkLogDetail(c *fiber.Ctx) error {
	return render(c, pages.WorkLogDetailPage())
}

func (h *Handler) CampaignDetail(c *fiber.Ctx) error {
	return render(c, pages.CampaignDetailPage())
}

func (h *Handler) KnowledgeDetail(c *fiber.Ctx) error {
	return render(c, pages.KnowledgeDetailPage())
}

func (h *Handler) InboxDetail(c *fiber.Ctx) error {
	return render(c, pages.InboxDetailPage())
}

func (h *Handler) InboxCreate(c *fiber.Ctx) error {
	return render(c, pages.InboxCreatePage())
}

func (h *Handler) BookmarkDetail(c *fiber.Ctx) error {
	return render(c, pages.BookmarkDetailPage())
}
