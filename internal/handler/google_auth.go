package handler

import (
	"crypto/rand"
	"encoding/hex"
	"strings"
	"time"

	"binhvuongos/internal/middleware"
	"binhvuongos/internal/oauth"

	"github.com/gofiber/fiber/v2"
)

// stateCookieName stores the CSRF state between /auth/google and /auth/google/callback.
const stateCookieName = "gstate"

// sessionTTL mirrors the non-remember JWT window used by email/password login.
const googleSessionTTL = 24 * time.Hour

func (h *Handler) newGoogleOAuth() *oauth.Config {
	return &oauth.Config{
		ClientID:     h.config.GoogleClientID,
		ClientSecret: h.config.GoogleClientSecret,
		RedirectURI:  h.config.GoogleRedirectURI,
	}
}

// GoogleLoginRedirect generates a CSRF state cookie and redirects to Google consent.
func (h *Handler) GoogleLoginRedirect(c *fiber.Ctx) error {
	if h.config.GoogleClientID == "" || h.config.GoogleClientSecret == "" {
		return c.Status(503).SendString("Google OAuth chưa được cấu hình")
	}
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return c.Status(500).SendString("Lỗi tạo state")
	}
	state := hex.EncodeToString(buf)
	c.Cookie(&fiber.Cookie{
		Name:     stateCookieName,
		Value:    state,
		Path:     "/",
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Lax",
		Expires:  time.Now().Add(10 * time.Minute),
	})
	return c.Redirect(h.newGoogleOAuth().AuthorizeURL(state))
}

// GoogleCallback exchanges the auth code, matches email against the users table,
// and issues a JWT cookie on success. Unknown or inactive emails are rejected —
// we never auto-provision accounts (admin-gated only).
func (h *Handler) GoogleCallback(c *fiber.Ctx) error {
	state := c.Query("state")
	cookieState := c.Cookies(stateCookieName)
	if state == "" || state != cookieState {
		return c.Status(400).SendString("State không hợp lệ (link có thể đã hết hạn, thử lại)")
	}
	c.ClearCookie(stateCookieName)

	code := c.Query("code")
	if code == "" {
		return c.Status(400).SendString("Thiếu code từ Google")
	}

	oc := h.newGoogleOAuth()
	accessToken, err := oc.ExchangeCode(c.Context(), code)
	if err != nil {
		return c.Status(502).SendString("Google exchange fail: " + err.Error())
	}
	info, err := oc.FetchUserinfo(c.Context(), accessToken)
	if err != nil || info.Email == "" {
		return c.Status(502).SendString("Không lấy được email từ Google")
	}
	if !info.EmailVerified {
		return c.Status(403).SendString("Email Google chưa được verify")
	}
	email := strings.ToLower(strings.TrimSpace(info.Email))

	user, err := h.queries.GetUserByEmail(c.Context(), email)
	if err != nil {
		return c.Status(403).SendString("Tài khoản " + email + " chưa được cấp quyền truy cập Bình Vương OS. Liên hệ admin.")
	}
	if user.Status != "active" {
		return c.Status(403).SendString("Tài khoản đã bị khoá")
	}

	// Populate avatar from Google picture if user hasn't set one yet.
	if (!user.AvatarURL.Valid || user.AvatarURL.String == "") && info.Picture != "" {
		_ = h.queries.UpdateUserAvatar(c.Context(), user.ID, info.Picture)
	}
	_ = h.queries.UpdateLastLogin(c.Context(), user.ID)

	jwtToken, err := middleware.GenerateToken(user, h.config.JWTSecret, false)
	if err != nil {
		return c.Status(500).SendString("Không tạo được JWT")
	}
	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    jwtToken,
		Path:     "/",
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Lax",
		Expires:  time.Now().Add(googleSessionTTL),
	})
	return c.Redirect("/inbox")
}
