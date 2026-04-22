package handler

import (
	"io"
	"regexp"

	"binhvuongos/internal/drive"

	"github.com/gofiber/fiber/v2"
)

// driveFileIDPattern whitelists Drive file ID characters to prevent path-injection.
// Google Drive file IDs are URL-safe base64-ish: alphanumerics, dash, underscore.
var driveFileIDPattern = regexp.MustCompile(`^[A-Za-z0-9_-]{10,64}$`)

// DriveProxy streams a Google Drive file through our server, reusing the
// app's refresh_token for Drive auth. Used to serve logos/avatars/inbox
// attachments as direct <img src=…> targets without making Drive files public.
//
// Auth: route is behind AuthRequired, so only logged-in users can fetch.
// Caching: Cache-Control private,max-age=3600 — per-browser cache for 1h.
func (h *Handler) DriveProxy(c *fiber.Ctx) error {
	fileID := c.Params("file_id")
	if !driveFileIDPattern.MatchString(fileID) {
		return c.Status(400).SendString("Invalid file id")
	}
	if h.config.GoogleRefreshToken == "" {
		return c.Status(503).SendString("Drive chưa được cấu hình")
	}
	cfg := &drive.Config{
		ClientID:     h.config.GoogleClientID,
		ClientSecret: h.config.GoogleClientSecret,
		RefreshToken: h.config.GoogleRefreshToken,
		FolderID:     h.config.GoogleDriveFolderID,
	}
	body, mime, err := drive.DownloadFile(c.Context(), cfg, fileID)
	if err != nil {
		return c.Status(502).SendString("Drive fetch fail: " + err.Error())
	}
	defer body.Close()

	c.Set("Content-Type", mime)
	c.Set("Cache-Control", "private, max-age=3600")
	if _, err := io.Copy(c.Response().BodyWriter(), body); err != nil {
		return c.Status(500).SendString("Stream fail")
	}
	return nil
}
