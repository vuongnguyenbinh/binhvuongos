package handler

import (
	"time"

	"binhvuongos/internal/middleware"
	"binhvuongos/web/templates/pages"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func (h *Handler) LoginPage(c *fiber.Ctx) error {
	// If already authenticated, redirect to dashboard
	if token := c.Cookies("token"); token != "" {
		return c.Redirect("/")
	}
	return render(c, pages.LoginPage(""))
}

func (h *Handler) Login(c *fiber.Ctx) error {
	email := c.FormValue("email")
	password := c.FormValue("password")
	rememberMe := c.FormValue("remember") == "on"

	if email == "" || password == "" {
		return render(c, pages.LoginPage("Vui lòng nhập email và mật khẩu"))
	}

	user, err := h.queries.GetUserByEmail(c.Context(), email)
	if err != nil {
		return render(c, pages.LoginPage("Email hoặc mật khẩu không đúng"))
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return render(c, pages.LoginPage("Email hoặc mật khẩu không đúng"))
	}

	token, err := middleware.GenerateToken(user, h.config.JWTSecret, rememberMe)
	if err != nil {
		return render(c, pages.LoginPage("Lỗi hệ thống, vui lòng thử lại"))
	}

	maxAge := 24 * 3600 // 1 day
	if rememberMe {
		maxAge = 30 * 24 * 3600 // 30 days
	}

	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    token,
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Lax",
		MaxAge:   maxAge,
		Path:     "/",
	})

	// Update last login
	_ = h.queries.UpdateLastLogin(c.Context(), user.ID)

	return c.Redirect("/")
}

func (h *Handler) Logout(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    "",
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Lax",
		MaxAge:   -1,
		Expires:  time.Now().Add(-1 * time.Hour),
		Path:     "/",
	})
	return c.Redirect("/login")
}

func (h *Handler) AuthMe(c *fiber.Ctx) error {
	user := GetUser(c)
	return c.JSON(fiber.Map{
		"id":        middleware.UUIDToString(user.ID),
		"email":     user.Email,
		"full_name": user.FullName,
		"role":      user.Role,
	})
}
