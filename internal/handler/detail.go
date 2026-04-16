package handler

import (
	"binhvuongos/web/templates/pages"
	"github.com/gofiber/fiber/v2"
)

func ContentDetail(c *fiber.Ctx) error {
	return render(c, pages.ContentDetailPage())
}

func CompanyDetail(c *fiber.Ctx) error {
	return render(c, pages.CompanyDetailPage())
}

func TaskDetail(c *fiber.Ctx) error {
	return render(c, pages.TaskDetailPage())
}

func WorkLogDetail(c *fiber.Ctx) error {
	return render(c, pages.WorkLogDetailPage())
}

func CampaignDetail(c *fiber.Ctx) error {
	return render(c, pages.CampaignDetailPage())
}

func KnowledgeDetail(c *fiber.Ctx) error {
	return render(c, pages.KnowledgeDetailPage())
}

func InboxDetail(c *fiber.Ctx) error {
	return render(c, pages.InboxDetailPage())
}

func InboxCreate(c *fiber.Ctx) error {
	return render(c, pages.InboxCreatePage())
}

func BookmarkDetail(c *fiber.Ctx) error {
	return render(c, pages.BookmarkDetailPage())
}
