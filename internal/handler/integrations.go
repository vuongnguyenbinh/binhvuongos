package handler

import (
	"database/sql"
	"log"

	"binhvuongos/internal/db/generated"

	"github.com/gofiber/fiber/v2"
)

// NotionSyncStatus returns sync status (stub - ready for Notion API integration)
func (h *Handler) NotionSyncStatus(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"status":       "not_configured",
			"message":      "Notion sync chưa được cấu hình. Cần NOTION_API_KEY trong .env",
			"last_sync":    nil,
			"tables_ready": []string{"users", "companies", "tasks", "content", "campaigns", "work_logs", "knowledge_items"},
		},
	})
}

// NotionSyncTrigger triggers a manual sync (stub)
func (h *Handler) NotionSyncTrigger(c *fiber.Ctx) error {
	// TODO: Implement Notion API sync
	// 1. Read NOTION_API_KEY from config
	// 2. For each table with sync_status='pending', push to Notion
	// 3. Log results to notion_sync_log table
	return c.JSON(fiber.Map{
		"success": false,
		"error":   "Notion sync chưa được triển khai. Cấu hình NOTION_API_KEY để bắt đầu.",
	})
}

// TelegramWebhook handles incoming Telegram bot messages (stub)
func (h *Handler) TelegramWebhook(c *fiber.Ctx) error {
	// TODO: Implement Telegram bot webhook
	// 1. Parse Update from Telegram
	// 2. Extract message text/photos
	// 3. Create inbox_items with source='telegram'
	// 4. Reply with confirmation

	var body struct {
		Message struct {
			Text string `json:"text"`
			From struct {
				ID       int64  `json:"id"`
				Username string `json:"username"`
			} `json:"from"`
			Chat struct {
				ID int64 `json:"id"`
			} `json:"chat"`
		} `json:"message"`
	}

	if err := c.BodyParser(&body); err != nil {
		log.Printf("Telegram webhook parse error: %v", err)
		return c.SendStatus(200) // Always return 200 to Telegram
	}

	if body.Message.Text != "" {
		// Auto-create inbox item from Telegram message
		_, err := h.queries.CreateInboxItem(c.Context(), generated.CreateInboxItemParams{
			Content:  body.Message.Text,
			Source:   sql.NullString{String: "telegram", Valid: true},
			ItemType: sql.NullString{String: "note", Valid: true},
		})
		if err != nil {
			log.Printf("Telegram inbox create error: %v", err)
		}
	}

	return c.SendStatus(200)
}
