package handler

import (
	"fmt"

	"binhvuongos/internal/drive"

	"github.com/gofiber/fiber/v2"
)

// Upload handles file upload to Google Drive and returns JSON with file URL
func (h *Handler) Upload(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "No file provided"})
	}

	// Max 50MB
	if file.Size > 50*1024*1024 {
		return c.Status(400).JSON(fiber.Map{"error": "File too large (max 50MB)"})
	}

	src, err := file.Open()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Cannot open file"})
	}
	defer src.Close()

	cfg := &drive.Config{
		ClientID:     h.config.GoogleClientID,
		ClientSecret: h.config.GoogleClientSecret,
		RefreshToken: h.config.GoogleRefreshToken,
		FolderID:     h.config.GoogleDriveFolderID,
	}

	// Skip if no Google credentials configured
	if cfg.RefreshToken == "" {
		return c.Status(500).JSON(fiber.Map{"error": "Google Drive not configured"})
	}

	result, err := drive.UploadFile(c.Context(), cfg, file.Filename, file.Header.Get("Content-Type"), src)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": fmt.Sprintf("Upload failed: %v", err)})
	}

	// Return Drive link
	driveURL := result.WebViewLink
	if driveURL == "" {
		driveURL = fmt.Sprintf("https://drive.google.com/file/d/%s/view", result.FileID)
	}

	// If HTMX request, return HTML snippet
	if c.Get("HX-Request") == "true" {
		return c.SendString(fmt.Sprintf(
			`<div class="flex items-center gap-2 text-xs text-forest mono"><span>✓</span><a href="%s" target="_blank" class="underline decoration-dotted">%s</a></div>`,
			driveURL, result.FileName))
	}

	return c.JSON(fiber.Map{
		"success":  true,
		"file_id":  result.FileID,
		"file_name": result.FileName,
		"url":      driveURL,
	})
}
