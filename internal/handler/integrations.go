package handler

import (
	"github.com/gofiber/fiber/v2"
)

// NotionSyncStatus returns sync status
func (h *Handler) NotionSyncStatus(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"status":       "not_configured",
			"message":      "Notion sync chưa được cấu hình. Cần NOTION_API_KEY trong .env",
			"tables_ready": []string{"users", "companies", "tasks", "content", "campaigns", "work_logs", "knowledge_items"},
		},
	})
}

// NotionSyncTrigger triggers a manual sync (stub)
func (h *Handler) NotionSyncTrigger(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": false,
		"error":   "Notion sync chưa được triển khai. Cấu hình NOTION_API_KEY để bắt đầu.",
	})
}
