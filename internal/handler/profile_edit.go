package handler

import (
	"strings"

	"binhvuongos/internal/db/generated"
	"binhvuongos/internal/drive"

	"github.com/gofiber/fiber/v2"
)

// maxAvatarBytes caps avatar uploads so users can't fill Drive with huge images.
const maxAvatarBytes = 10 * 1024 * 1024

// ProfileUpdate lets a user edit their own full_name + phone.
// Email and role stay immutable (managed by admin).
func (h *Handler) ProfileUpdate(c *fiber.Ctx) error {
	user := GetUser(c)
	fullName := strings.TrimSpace(c.FormValue("full_name"))
	if fullName == "" {
		return c.Status(400).SendString("Họ tên không được để trống")
	}
	if err := h.queries.UpdateOwnProfile(c.Context(), generated.UpdateOwnProfileParams{
		ID:       user.ID,
		FullName: fullName,
		Phone:    strings.TrimSpace(c.FormValue("phone")),
	}); err != nil {
		return c.Status(500).SendString("Lỗi cập nhật")
	}
	return c.Redirect("/profile")
}

// ProfileAvatar uploads an image to Drive and writes the resulting URL
// to users.avatar_url. Reuses the existing drive.UploadFile helper.
func (h *Handler) ProfileAvatar(c *fiber.Ctx) error {
	user := GetUser(c)
	file, err := c.FormFile("avatar")
	if err != nil {
		return c.Status(400).SendString("Thiếu file avatar")
	}
	if file.Size > maxAvatarBytes {
		return c.Status(400).SendString("Ảnh quá lớn (tối đa 10MB)")
	}
	if h.config.GoogleRefreshToken == "" {
		return c.Status(503).SendString("Drive chưa được cấu hình")
	}
	src, err := file.Open()
	if err != nil {
		return c.Status(500).SendString("Không mở được file")
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
		return c.Status(500).SendString("Upload fail: " + err.Error())
	}
	// Store internal proxy path so <img src> renders through our server.
	avatarURL := "/drive/" + result.FileID
	if err := h.queries.UpdateUserAvatar(c.Context(), user.ID, avatarURL); err != nil {
		return c.Status(500).SendString("Lỗi lưu URL")
	}
	return c.Redirect("/profile")
}
