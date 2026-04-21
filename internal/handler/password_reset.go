package handler

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"binhvuongos/internal/middleware"
	"binhvuongos/web/templates/pages"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

const minPasswordLen = 8

// GenerateResetLink creates a 1-hour single-use reset token and renders the URL
// for the actor to share manually. SMTP-ready: drop a mailer.Send() call here later.
func (h *Handler) GenerateResetLink(c *fiber.Ctx) error {
	actor := GetUser(c)
	target, err := h.queries.GetUserByID(c.Context(), middleware.StringToUUID(c.Params("id")))
	if err != nil {
		return c.Status(404).SendString("Không tìm thấy user")
	}
	if !middleware.CanManageUser(actor, target) {
		return c.Status(403).SendString("Không có quyền reset password user này")
	}
	buf := make([]byte, 16) // 128-bit entropy
	if _, err := rand.Read(buf); err != nil {
		return c.Status(500).SendString("Lỗi tạo token")
	}
	token := hex.EncodeToString(buf)
	if _, err := h.queries.CreateResetToken(c.Context(), target.ID, token); err != nil {
		return c.Status(500).SendString("Lỗi lưu token")
	}
	resetURL := fmt.Sprintf("https://os.binhvuong.vn/reset/%s", token)
	// TODO: when SMTP added, h.mailer.SendResetEmail(target.Email, resetURL)
	return render(c, pages.ResetLinkFlash(target.Email, resetURL))
}

// ResetPasswordPage renders the public form when the token is still valid.
func (h *Handler) ResetPasswordPage(c *fiber.Ctx) error {
	token := c.Params("token")
	if _, err := h.queries.GetValidResetToken(c.Context(), token); err != nil {
		return c.Status(404).SendString("Link đặt lại mật khẩu không hợp lệ hoặc đã hết hạn")
	}
	return render(c, pages.PasswordResetForm(token))
}

// ResetPassword validates the token, sets the new password, marks the token used.
func (h *Handler) ResetPassword(c *fiber.Ctx) error {
	token := c.Params("token")
	row, err := h.queries.GetValidResetToken(c.Context(), token)
	if err != nil {
		return c.Status(404).SendString("Link đặt lại mật khẩu không hợp lệ hoặc đã hết hạn")
	}
	password := c.FormValue("password")
	if len(password) < minPasswordLen {
		return c.Status(400).SendString(fmt.Sprintf("Mật khẩu cần ít nhất %d ký tự", minPasswordLen))
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return c.Status(500).SendString("Lỗi mã hoá mật khẩu")
	}
	if err := h.queries.UpdatePassword(c.Context(), row.UserID, string(hash)); err != nil {
		return c.Status(500).SendString("Lỗi cập nhật mật khẩu")
	}
	_ = h.queries.MarkResetTokenUsed(c.Context(), token)
	return c.Redirect("/login?reset=1")
}
