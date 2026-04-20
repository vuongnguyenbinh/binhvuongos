package handler

import (
	"errors"
	"fmt"
	"mime/multipart"
	"strings"

	"binhvuongos/internal/drive"

	"github.com/gofiber/fiber/v2"
)

// allowedItemTypes is the closed whitelist for inbox item_type.
// Other values are rejected by validateInboxInput to keep DB tidy.
var allowedItemTypes = map[string]bool{
	"note":  true,
	"link":  true,
	"image": true,
	"voice": true,
	"file":  true,
}

const (
	maxAttachments    = 10
	maxContentBytes   = 10 * 1024        // 10KB per message
	maxMultipartBytes = 50 * 1024 * 1024 // 50MB per uploaded file
	maxAttachmentURL  = 2048
	maxSourceLen      = 30
	maxExternalRefLen = 200
)

// detectItemType infers item_type when client does not supply one.
// Priority: attachment mime (image/audio/file) → URL present → note.
func detectItemType(content, url string, attachments []map[string]any) string {
	for _, a := range attachments {
		mime, _ := a["mime"].(string)
		switch {
		case strings.HasPrefix(mime, "image/"):
			return "image"
		case strings.HasPrefix(mime, "audio/"):
			return "voice"
		case mime != "":
			return "file"
		}
	}
	if url != "" || isHTTPURL(content) {
		return "link"
	}
	return "note"
}

// isHTTPURL is a cheap prefix check — good enough for inbox triage.
func isHTTPURL(s string) bool {
	return len(s) >= 7 && (s[:7] == "http://" || (len(s) >= 8 && s[:8] == "https://"))
}

// validateInboxInput enforces size/whitelist constraints shared by JSON + multipart paths.
func validateInboxInput(content, source, itemType, externalRef string, attachmentURLs []string) error {
	if content == "" {
		return errors.New("content is required")
	}
	if len(content) > maxContentBytes {
		return fmt.Errorf("content too large (max %d bytes)", maxContentBytes)
	}
	if source == "" {
		return errors.New("source is required")
	}
	if len(source) > maxSourceLen {
		return fmt.Errorf("source too long (max %d chars)", maxSourceLen)
	}
	if itemType != "" && !allowedItemTypes[itemType] {
		return errors.New("invalid item_type (allowed: note, link, image, voice, file)")
	}
	if len(externalRef) > maxExternalRefLen {
		return fmt.Errorf("external_ref too long (max %d chars)", maxExternalRefLen)
	}
	if len(attachmentURLs) > maxAttachments {
		return fmt.Errorf("too many attachments (max %d)", maxAttachments)
	}
	for _, u := range attachmentURLs {
		if len(u) > maxAttachmentURL {
			return fmt.Errorf("attachment URL too long (max %d chars)", maxAttachmentURL)
		}
	}
	return nil
}

// uploadAttachmentToDrive streams a multipart file to Google Drive and returns a structured attachment entry.
// Returns error if Drive is not configured or file exceeds size cap.
func (h *Handler) uploadAttachmentToDrive(c *fiber.Ctx, file *multipart.FileHeader) (map[string]any, error) {
	if file.Size > maxMultipartBytes {
		return nil, fmt.Errorf("file too large (max %d bytes)", maxMultipartBytes)
	}
	if h.config.GoogleRefreshToken == "" {
		return nil, errors.New("Google Drive not configured")
	}
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("cannot open file: %w", err)
	}
	defer src.Close()

	cfg := &drive.Config{
		ClientID:     h.config.GoogleClientID,
		ClientSecret: h.config.GoogleClientSecret,
		RefreshToken: h.config.GoogleRefreshToken,
		FolderID:     h.config.GoogleDriveFolderID,
	}
	result, err := drive.UploadFile(c.Context(), cfg, file.Filename, file.Header.Get("Content-Type"), src)
	if err != nil {
		return nil, fmt.Errorf("drive upload: %w", err)
	}
	url := result.WebViewLink
	if url == "" {
		url = fmt.Sprintf("https://drive.google.com/file/d/%s/view", result.FileID)
	}
	return map[string]any{
		"url":           url,
		"filename":      result.FileName,
		"mime":          file.Header.Get("Content-Type"),
		"drive_file_id": result.FileID,
	}, nil
}
