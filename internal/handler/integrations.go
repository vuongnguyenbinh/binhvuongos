package handler

import (
	"database/sql"
	"fmt"
	"log"

	"binhvuongos/internal/db/generated"

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

// TelegramWebhook handles incoming Telegram bot messages
func (h *Handler) TelegramWebhook(c *fiber.Ctx) error {
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
		return c.SendStatus(200)
	}

	// Verify telegram user exists in DB
	if body.Message.From.ID == 0 {
		return c.SendStatus(200)
	}

	telegramID := fmt.Sprintf("%d", body.Message.From.ID)

	// Check if telegram_id matches a user (basic verification)
	user, err := h.queries.GetUserByTelegramID(c.Context(), telegramID)
	if err != nil {
		log.Printf("Telegram: unknown user telegram_id=%s", telegramID)
		return c.SendStatus(200)
	}

	if body.Message.Text != "" {
		// Detect URL in message
		var url sql.NullString
		var itemType sql.NullString
		if isURL(body.Message.Text) {
			url = sql.NullString{String: body.Message.Text, Valid: true}
			itemType = sql.NullString{String: "link", Valid: true}
		} else {
			itemType = sql.NullString{String: "note", Valid: true}
		}

		_, err := h.queries.CreateInboxItem(c.Context(), generated.CreateInboxItemParams{
			Content:     body.Message.Text,
			URL:         url,
			Source:      sql.NullString{String: "telegram", Valid: true},
			ItemType:    itemType,
			SubmittedBy: user.ID,
		})
		if err != nil {
			log.Printf("Telegram inbox create error: %v", err)
		}
	}

	return c.SendStatus(200)
}

func isURL(s string) bool {
	return len(s) > 8 && (s[:7] == "http://" || s[:8] == "https://")
}
